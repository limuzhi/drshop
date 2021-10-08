package biz

import (
	"context"
	v1 "drpshop/api/task/v1"
	"drpshop/internal/apps/sys/model"
	"drpshop/pkg/httpclient"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jakecoffman/cron"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskServiceUsecase struct {
	taskRepo     SysTaskRepo
	taskGrpcRepo TaskGrpcPoolRepo
	log          *log.Helper

	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 任务计数-正在运行的任务
	taskCount *TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue *ConcurrencyQueue

	// 同一任务是否有实例处于运行中
	runInstance *Instance
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

func NewTaskServiceUsecase(repo SysTaskRepo, taskGrpcRepo TaskGrpcPoolRepo, logger log.Logger) *TaskServiceUsecase {
	return &TaskServiceUsecase{
		taskRepo:         repo,
		taskGrpcRepo:     taskGrpcRepo,
		log:              log.NewHelper(log.With(logger, "module", "sys/biz/task_servie")),
		serviceCron:      cron.New(),
		concurrencyQueue: &ConcurrencyQueue{queue: make(chan struct{}, 500)},
		runInstance:      &Instance{m: sync.Map{}},
		taskCount: &TaskCount{
			wg:   sync.WaitGroup{},
			exit: make(chan struct{}),
		},
	}
}

func (ts *TaskServiceUsecase) Initialize() {
	ts.serviceCron.Start()
	go ts.taskCount.Wait()
	ts.log.Info("开始初始化定时任务")
	taskNum := 0
	page := 1
	pageSize := 1000
	for {
		taskList, err := ts.taskRepo.ActiveList(context.Background(), page, pageSize)
		if err != nil {
			ts.log.Fatalf("定时任务初始化#获取任务列表错误: %s", err)
			break
		}
		if len(taskList) == 0 {
			break
		}
		for _, item := range taskList {
			ts.Add(item)
			taskNum++
		}
		page++
	}
	ts.log.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// 添加任务
func (ts *TaskServiceUsecase) Add(in *model.SysTask) {
	if in.Level == model.TaskLevelChild {
		ts.log.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", in.TaskId)
		return
	}
	taskFunc := ts.createJob(in)
	if taskFunc == nil {
		ts.log.Error("创建任务处理Job失败,不支持的任务协议#", in.Protocol)
		return
	}

	cronName := strconv.Itoa(int(in.TaskId))
	err := panicToError(func() {
		ts.serviceCron.AddFunc(in.Spec, taskFunc, cronName)
	})
	if err != nil {
		ts.log.Error("添加任务到调度器失败#", err)
	}
}

// 批量添加任务
func (ts *TaskServiceUsecase) BatchAdd(ins []*model.SysTask) {
	for _, item := range ins {
		ts.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (ts *TaskServiceUsecase) RemoveAndAdd(in *model.SysTask) {
	ts.Remove(in.TaskId)
	ts.Add(in)
}

func (ts *TaskServiceUsecase) NextRunTime(in *model.SysTask) time.Time {
	if in.Level != model.TaskLevelParent || in.Status != model.Enabled {
		return time.Time{}
	}
	entries := ts.serviceCron.Entries()
	taskName := strconv.Itoa(int(in.TaskId))
	for _, item := range entries {
		if item.Name == taskName {
			return item.Next
		}
	}
	return time.Time{}
}

// 停止运行中的任务
func (ts *TaskServiceUsecase) Stop(ip string, port int, id int64) {
	ts.taskGrpcRepo.Stop(ip, port, id)
}

//删除
func (ts *TaskServiceUsecase) Remove(taskId int64) {
	ts.serviceCron.RemoveJob(strconv.Itoa(int(taskId)))
}

//等待所有任务结束后退出
func (ts *TaskServiceUsecase) WaitAndExit() {
	ts.serviceCron.Stop()
	ts.taskCount.Exit()
}

func (ts *TaskServiceUsecase) Run(in *model.SysTask) {
	go ts.createJob(in)()
}

// 直接运行任务

func panicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			stackBuf := make([]byte, 4096)
			n := runtime.Stack(stackBuf, false)
			err = fmt.Errorf(fmt.Sprintf("panic: %v %s", err, stackBuf[:n]))
		}
	}()
	f()
	return
}

func (ts *TaskServiceUsecase) createJob(in *model.SysTask) cron.FuncJob {
	handler := ts.createHandler(in.Protocol)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		ts.taskCount.Add()
		defer ts.taskCount.Done()

		taskLogId := ts.beforeExecJob(in)
		if taskLogId <= 0 {
			return
		}
		if in.Multi == 0 {
			ts.runInstance.add(in.TaskId)
			defer ts.runInstance.done(in.TaskId)
		}
		ts.concurrencyQueue.Add()
		defer ts.concurrencyQueue.Done()

		ts.log.Infof("开始执行任务#%s#命令-%s", in.Name, in.Command)
		taskResult := ts.execJob(handler, in, taskLogId)
		ts.log.Infof("任务完成#%s#命令-%s", in.Name, in.Command)
		ts.afterExecJob(in, taskResult, taskLogId)
	}
	return taskFunc
}

func (ts *TaskServiceUsecase) beforeExecJob(in *model.SysTask) int64 {
	if in.Multi == 0 || ts.runInstance.has(in.TaskId) {
		ts.createTaskLog(in, model.Cancel)
		return 0
	}
	taskLogId, err := ts.createTaskLog(in, model.Running)
	if err != nil {
		ts.log.Error("任务开始执行#写入任务日志失败-", err)
		return 0
	}
	ts.log.Debugf("任务命令-%s", in.Command)
	return taskLogId
}

//执行具体任务
func (ts *TaskServiceUsecase) execJob(handler Handler, in *model.SysTask, taskLogId int64) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			ts.log.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if in.RetryTimes > 0 {
		execTimes += in.RetryTimes
	}

	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(in, taskLogId, ts.taskGrpcRepo)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			ts.log.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", in.TaskId, i, output, err.Error())
			if in.RetryInterval > 0 {
				time.Sleep(time.Duration(in.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}
	return TaskResult{Result: output, Err: err, RetryTimes: in.RetryTimes}
}

// 任务执行后置操作
func (ts *TaskServiceUsecase) afterExecJob(in *model.SysTask, taskResult TaskResult, taskLogId int64) {

	if err := ts.updateTaskLog(taskLogId, taskResult); err != nil {
		ts.log.Error("任务结束#更新任务日志失败-", err)
	}
	// 发送邮件
	go ts.SendNotification(in, taskResult)

	// 执行依赖任务
	go ts.execDependencyTask(in, taskResult)
}

// 执行依赖任务, 多个任务并发执行
func (ts *TaskServiceUsecase) execDependencyTask(in *model.SysTask, taskResult TaskResult) {
	// 父任务才能执行子任务
	if in.Level != model.TaskLevelParent {
		return
	}
	// 是否存在子任务
	dependencyTaskId := strings.TrimSpace(in.DependencyTaskId)
	if dependencyTaskId == "" {
		return
	}

	// 父子任务关系为强依赖, 父任务执行失败, 不执行依赖任务
	if in.DependencyStatus == model.TaskDependencyStatusStrong && taskResult.Err != nil {
		ts.log.Infof("父子任务为强依赖关系, 父任务执行失败, 不运行依赖任务#主任务ID-%d", in.TaskId)
		return
	}

	//获取子任务
	tasks, err := ts.taskRepo.GetDependencyTaskList(context.Background(), dependencyTaskId)
	if err != nil {
		ts.log.Errorf("获取依赖任务失败#主任务ID-%d#%s", in.TaskId, err.Error())
		return
	}
	if len(tasks) == 0 {
		ts.log.Errorf("依赖任务列表为空#主任务ID-%d", in.TaskId)
	}
	for _, task := range tasks {
		task.Spec = fmt.Sprintf("依赖任务(主任务ID-%d)", in.TaskId)
		ts.Run(task)
	}
}

// 创建任务日志
func (ts *TaskServiceUsecase) createTaskLog(in *model.SysTask, status model.Status) (int64, error) {
	//TODO
	taskLogInsert := &model.SysTaskLog{
		TaskId:   in.TaskId,
		Name:     in.Name,
		Spec:     in.Spec,
		Protocol: in.Protocol,
		Command:  in.Command,
		Timeout:  in.Timeout,
	}
	if in.Protocol == model.TaskRPC {
		aggregationHost := ""
		for _, host := range in.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.Alias, host.Name)
		}
		taskLogInsert.Hostname = aggregationHost
	}
	taskLogInsert.StartTime = time.Now().Unix()
	taskLogInsert.Status = status
	//TODO insert
	ts.log.Info("TODO")
	return taskLogInsert.TaskId, nil
}

// 更新任务日志
func (ts *TaskServiceUsecase) updateTaskLog(taskLogId int64, taskResult TaskResult) error {
	var status model.Status
	if taskResult.Err != nil {
		status = model.Failure
	} else {
		status = model.Finish
	}
	taskLogMap := make(map[string]interface{})
	taskLogMap["retry_times"] = taskResult.RetryTimes
	taskLogMap["status"] = status
	taskLogMap["result"] = taskResult.Result
	//TODO update
	ts.log.Info("TODO")
	return nil
}

func (ts *TaskServiceUsecase) createHandler(protocol model.TaskProtocol) Handler {
	var handler Handler = nil
	switch protocol {
	case model.TaskHTTP:
		handler = new(HTTPHandler)
	case model.TaskRPC:
		handler = new(RPCHandler)
	}
	return handler
}

// 发送任务结果通知
func (ts *TaskServiceUsecase) SendNotification(in *model.SysTask, taskResult TaskResult) {
	var statusName string
	// 未开启通知
	if in.NotifyStatus == 0 {
		return
	}
	if in.NotifyStatus == 3 {
		// 关键字匹配通知
		if !strings.Contains(taskResult.Result, in.NotifyKeyword) {
			return
		}
	}
	if in.NotifyStatus == 1 && taskResult.Err == nil {
		// 执行失败才发送通知
		return
	}
	if in.NotifyType != 3 && in.NotifyReceiverId == "" {
		return
	}
	if taskResult.Err != nil {
		statusName = "失败"
	} else {
		statusName = "成功"
	}
	// 发送通知
	msg := Message{
		"task_type":        in.NotifyType,
		"task_receiver_id": in.NotifyReceiverId,
		"name":             in.Name,
		"output":           taskResult.Result,
		"status":           statusName,
		"task_id":          in.TaskId,
		"remark":           in.Remark,
	}
	Push(msg)
	fmt.Println(statusName)
}

type Handler interface {
	Run(in *model.SysTask, taskId int64, taskGrpcRepo TaskGrpcPoolRepo) (string, error)
}

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(in *model.SysTask, taskId int64, taskGrpcRepo TaskGrpcPoolRepo) (result string, err error) {
	if in.Timeout <= 0 || in.Timeout > HttpExecTimeout {
		in.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if in.HttpMethod == model.TaskHTTPMethodGet {
		resp = httpclient.Get(in.Command, in.Timeout)
	} else {
		urlFields := strings.Split(in.Command, "?")
		in.Command = urlFields[0]
		var params string
		for len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(in.Command, params, in.Timeout)
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}
	return resp.Body, err
}

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(in *model.SysTask, taskId int64, taskGrpcRepo TaskGrpcPoolRepo) (result string, err error) {
	taskReq := &v1.TaskReq{
		Timeout: int32(in.Timeout),
		Command: in.Command,
		Id:      in.TaskId,
	}
	resultChan := make(chan TaskResult, len(in.Hosts))
	for _, taskHost := range in.Hosts {
		go func(th *model.SysHost) {
			output, err := taskGrpcRepo.Exec(th.Name, th.Port, taskReq)
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s:%d]\n%s\n%s\n\n",
				th.Alias, th.Name, th.Port, errorMessage, output,
			)
			resultChan <- TaskResult{Err: err, Result: outputMessage}
		}(taskHost)
	}
	var aggregationErr error = nil
	aggregationResult := ""
	for i := 0; i < len(in.Hosts); i++ {
		taskResult := <-resultChan
		aggregationResult += taskResult.Result
		if taskResult.Err != nil {
			aggregationErr = taskResult.Err
		}
	}
	return aggregationResult, aggregationErr
}

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key int64) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key int64) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key int64) {
	i.m.Delete(key)
}

// 任务计数
type TaskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *TaskCount) Add() {
	tc.wg.Add(1)
}

func (tc *TaskCount) Done() {
	tc.wg.Done()
}

func (tc *TaskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *TaskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}
