package service

import (
	"context"
	v1 "drpshop/api/tasknode/v1"
	"drpshop/internal/apps/sys/biz"
	"drpshop/internal/apps/sys/service/notify"
	"drpshop/internal/apps/sys/service/taskrpc"
	"drpshop/pkg/httpclient"
	"errors"
	"fmt"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"drpshop/internal/apps/sys/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jakecoffman/cron"
	"github.com/ouqiang/goutil"
)

type TaskService struct {
	repo biz.SysTaskRepo
	log  *log.Helper
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

func NewTaskService(repo biz.SysTaskRepo, logger log.Logger) *TaskService {
	return &TaskService{repo: repo, log: log.NewHelper(
		log.With(logger, "module", "sys/service/task_service"))}
}

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 同一任务是否有实例处于运行中
	runInstance Instance

	// 任务计数-正在运行的任务
	taskCount TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue ConcurrencyQueue

	taskMap sync.Map
)

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key int) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key int) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key int) {
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

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task TaskService) Initialize() {
	serviceCron = cron.New()
	serviceCron.Start()
	concurrencyQueue = ConcurrencyQueue{
		queue: make(chan struct{}, 500),
	}
	taskCount = TaskCount{
		sync.WaitGroup{},
		make(chan struct{}),
	}
	go taskCount.Wait()

	task.log.Info("开始初始化定时任务")
	taskNum := 0
	page := 1
	pageSize := 1000

	for {
		taskList, err := task.repo.ActiveList(context.Background(), page, pageSize)
		if err != nil {
			task.log.Fatalf("定时任务初始化#获取任务列表错误: %s", err)
			break
		}
		if len(taskList) == 0 {
			break
		}
		for _, item := range taskList {
			task.Add(item)
			taskNum++
		}
		page++
	}
	task.log.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// 批量添加任务
func (task TaskService) BatchAdd(list []*model.SysTask) {
	for _, item := range list {
		task.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (task TaskService) RemoveAndAdd(taskModel *model.SysTask) {
	task.Remove(int(taskModel.TaskId))
	task.Add(taskModel)
}

//// 添加任务
func (task TaskService) Add(taskModel *model.SysTask) {
	if taskModel.Level == model.TaskLevelChild {
		task.log.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.TaskId)
		return
	}
	taskFunc := task.createJob(taskModel)
	if taskFunc == nil {
		task.log.Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return
	}

	cronName := strconv.Itoa(int(taskModel.TaskId))
	err := goutil.PanicToError(func() {
		serviceCron.AddFunc(taskModel.Spec, taskFunc, cronName)
	})
	if err != nil {
		task.log.Error("添加任务到调度器失败#", err)
	}
}

// 停止运行中的任务
func (task TaskService) Stop(ip string, port int, id int64) {
	//rpcClient.Stop(ip, port, id)
}

func (task TaskService) Remove(id int) {
	serviceCron.RemoveJob(strconv.Itoa(id))
}

// 等待所有任务结束后退出
func (task TaskService) WaitAndExit() {
	serviceCron.Stop()
	taskCount.Exit()
}

// 直接运行任务
func (task TaskService) Run(taskModel *model.SysTask) {
	go task.createJob(taskModel)()
}

func (task TaskService) NextRunTime(taskModel *model.SysTask) time.Time {
	if taskModel.Level != model.TaskLevelParent ||
		taskModel.Status != model.Enabled {
		return time.Time{}
	}
	entries := serviceCron.Entries()
	taskName := strconv.Itoa(int(taskModel.TaskId))
	for _, item := range entries {
		if item.Name == taskName {
			return item.Next
		}
	}
	return time.Time{}
}

func (task TaskService) createJob(taskModel *model.SysTask) cron.FuncJob {
	handler := createHandler(taskModel)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()

		taskLogId := task.beforeExecJob(taskModel)
		if taskLogId <= 0 {
			return
		}

		if taskModel.Multi == 0 {
			runInstance.add(int(taskModel.TaskId))
			defer runInstance.done(int(taskModel.TaskId))
		}

		concurrencyQueue.Add()
		defer concurrencyQueue.Done()

		task.log.Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
		taskResult := task.execJob(handler, taskModel, taskLogId)
		task.log.Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
		task.afterExecJob(taskModel, taskResult, taskLogId)
	}

	return taskFunc
}

// 任务前置操作
func (task TaskService) beforeExecJob(taskModel *model.SysTask) (taskLogId int64) {
	if taskModel.Multi == 0 && runInstance.has(int(taskModel.TaskId)) {
		task.createTaskLog(taskModel, model.Cancel)
		return
	}
	taskLogId, err := task.createTaskLog(taskModel, model.Running)
	if err != nil {
		task.log.Error("任务开始执行#写入任务日志失败-", err)
		return
	}
	task.log.Debugf("任务命令-%s", taskModel.Command)

	return taskLogId
}

// 任务执行后置操作
func (task TaskService) afterExecJob(taskModel *model.SysTask, taskResult TaskResult, taskLogId int64) {
	_, err := task.updateTaskLog(taskLogId, taskResult)
	if err != nil {
		task.log.Error("任务结束#更新任务日志失败-", err)
	}

	// 发送邮件
	go task.sendNotification(taskModel, taskResult)
	// 执行依赖任务
	go task.execDependencyTask(taskModel, taskResult)
}

// 创建任务日志
func (task TaskService) createTaskLog(taskModel *model.SysTask, status model.Status) (int64, error) {
	taskLogModel := &model.SysTaskLog{}
	taskLogModel.TaskId = taskModel.TaskId
	taskLogModel.Name = taskModel.Name
	taskLogModel.Spec = taskModel.Spec
	taskLogModel.Protocol = taskModel.Protocol
	taskLogModel.Command = taskModel.Command
	taskLogModel.Timeout = taskModel.Timeout
	if taskModel.Protocol == model.TaskRPC {
		aggregationHost := ""
		for _, host := range taskModel.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.Alias, host.Name)
		}
		taskLogModel.Hostname = aggregationHost
	}
	taskLogModel.StartTime = time.Now().Unix()
	taskLogModel.Status = status
	//TODO

	return 0, nil
}

// 更新任务日志
func (task TaskService) updateTaskLog(taskLogId int64, taskResult TaskResult) (int64, error) {
	taskLogModel := new(model.SysTaskLog)
	var status model.Status
	result := taskResult.Result
	if taskResult.Err != nil {
		status = model.Failure
	} else {
		status = model.Finish
	}
	fmt.Println(taskLogModel, result, status)
	//return taskLogModel.Update(taskLogId, models.CommonMap{
	//	"retry_times": taskResult.RetryTimes,
	//	"status":      status,
	//	"result":      result,
	//})
	return 0, nil
}

// 发送任务结果通知
func (task TaskService) sendNotification(taskModel *model.SysTask, taskResult TaskResult) {
	var statusName string
	// 未开启通知
	if taskModel.NotifyStatus == 0 {
		return
	}
	if taskModel.NotifyStatus == 3 {
		// 关键字匹配通知
		if !strings.Contains(taskResult.Result, taskModel.NotifyKeyword) {
			return
		}
	}
	if taskModel.NotifyStatus == 1 && taskResult.Err == nil {
		// 执行失败才发送通知
		return
	}
	if taskModel.NotifyType != 3 && taskModel.NotifyReceiverId == "" {
		return
	}
	if taskResult.Err != nil {
		statusName = "失败"
	} else {
		statusName = "成功"
	}
	// 发送通知
	msg := notify.Message{
		"task_type":        taskModel.NotifyType,
		"task_receiver_id": taskModel.NotifyReceiverId,
		"name":             taskModel.Name,
		"output":           taskResult.Result,
		"status":           statusName,
		"task_id":          taskModel.TaskId,
		"remark":           taskModel.Remark,
	}
	notify.Push(msg)
}

// 执行依赖任务, 多个任务并发执行
func (task TaskService) execDependencyTask(taskModel *model.SysTask, taskResult TaskResult) {
	// 父任务才能执行子任务
	if taskModel.Level != model.TaskLevelParent {
		return
	}
	// 是否存在子任务
	dependencyTaskId := strings.TrimSpace(taskModel.DependencyTaskId)
	if dependencyTaskId == "" {
		return
	}

	// 父子任务关系为强依赖, 父任务执行失败, 不执行依赖任务
	if taskModel.DependencyStatus == model.TaskDependencyStatusStrong && taskResult.Err != nil {
		task.log.Infof("父子任务为强依赖关系, 父任务执行失败, 不运行依赖任务#主任务ID-%d", taskModel.TaskId)
		return
	}

	// 获取子任务
	list, err := task.repo.GetDependencyTaskList(context.Background(), dependencyTaskId)
	if err != nil {
		task.log.Errorf("获取依赖任务失败#主任务ID-%d#%s", taskModel.TaskId, err.Error())
		return
	}
	if len(list) == 0 {
		task.log.Errorf("依赖任务列表为空#主任务ID-%d", taskModel.TaskId)
	}
	for _, v := range list {
		v.Spec = fmt.Sprintf("依赖任务(主任务ID-%d)", taskModel.TaskId)
		task.Run(v)
	}
}

// 执行具体任务
func (task TaskService) execJob(handler Handler, taskModel *model.SysTask, taskId int64) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			task.log.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if taskModel.RetryTimes > 0 {
		execTimes += taskModel.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(taskModel, taskId)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			task.log.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s",
				taskModel.TaskId, i, output, err.Error())
			if taskModel.RetryInterval > 0 {
				time.Sleep(time.Duration(taskModel.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}

	return TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}

func createHandler(taskModel *model.SysTask) Handler {
	var handler Handler = nil
	switch taskModel.Protocol {
	case model.TaskHTTP:
		handler = new(HTTPHandler)
	case model.TaskRPC:
		handler = new(RPCHandler)
	}

	return handler
}

type Handler interface {
	Run(taskModel *model.SysTask, taskId int64) (string, error)
}

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(taskModel *model.SysTask, taskId int64) (string, error) {
	if taskModel.Timeout <= 0 || taskModel.Timeout > HttpExecTimeout {
		taskModel.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if taskModel.HttpMethod == model.TaskHTTPMethodGet {
		resp = httpclient.Get(taskModel.Command, taskModel.Timeout)
	} else {
		urlFields := strings.Split(taskModel.Command, "?")
		taskModel.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(taskModel.Command, params, taskModel.Timeout)
	}
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}

	return resp.Body, nil
}

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel *model.SysTask, taskId int64) (string, error) {
	taskReq := new(v1.TaskReq)
	taskReq.Timeout = int32(taskModel.Timeout)
	taskReq.Command = taskModel.Command
	taskReq.Id = taskId
	resultChan := make(chan TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		go func(th *model.SysHost) {
			output, err := taskrpc.Exec(th.Name, th.Port, taskReq)
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
	for i := 0; i < len(taskModel.Hosts); i++ {
		taskResult := <-resultChan
		aggregationResult += taskResult.Result
		if taskResult.Err != nil {
			aggregationErr = taskResult.Err
		}
	}

	return aggregationResult, aggregationErr
}

func NewTasknodeServiceClient(r registry.Discovery) v1.TaskClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///beer.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := v1.NewTaskClient(conn)
	return c
}

func NewDiscovery(address, scheme string) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = address
	c.Scheme = scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errors.New("无法连接远程服务器")
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
