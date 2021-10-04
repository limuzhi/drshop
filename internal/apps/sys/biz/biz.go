package biz

import "github.com/google/wire"

//TODO
// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSysApisUsecase,NewSysUserUsecase,
	NewSysRoleUsecase,NewSysLoginLogUsecase,NewSysConfigUsecase,
	NewSysDictDataUsecase,NewSysDictTypeUsecase,NewSysJobUsecase,
	NewSysMenuUsecase,NewSysPostUsecase,NewSysHostUsecase,
	NewSysTaskUsecase,NewSysDeptUsecase)

