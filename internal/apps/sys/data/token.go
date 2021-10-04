package data

import "drpshop/pkg/token"

func (r *sysUserRepo) CreateToken(claims *token.UserClaims) (string, error) {
	return r.data.tk.CreateToken(claims)
}
