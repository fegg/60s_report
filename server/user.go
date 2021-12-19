package calendar

import (
	"context"
	"fmt"
	medCalendar "med-common/app/service/med-calendar/api/v1"
	"med-common/app/service/med-calendar/internal/pkg/toutils"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"
	"strings"

	"github.com/go-kratos/kratos/pkg/log"

	"gorm.io/gorm"
)

func (d *Dao) GetTopDeptByUserId(userId string) (*medCalendar.GetTopDeptByUserIdResp, error) {
	resp := &medCalendar.GetTopDeptByUserIdResp{
		List: make([]*medCalendar.GetTopDeptByUserIdResp_DeptInfo, 0),
	}
	userInfo := &model.User{}
	err := d.orm.Model(&model.User{}).First(userInfo, "user_id = ?", userId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, code.MedErrDBFailed
	}
	dept := &model.Dept{}
	err = d.orm.Model(&model.Dept{}).Where("dept_id = ?", userInfo.DeptId).First(dept).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, code.MedErrDBFailed
	}
	var level int64 = 1
	resp.List = append(resp.List, &medCalendar.GetTopDeptByUserIdResp_DeptInfo{
		Level: level,
		Dept: &medCalendar.GetTopDeptByUserIdResp_Dept{
			DeptId:   dept.DeptId,
			DeptName: dept.Name,
		},
	})
	if dept.Level == 1 {
		return resp, nil
	}
	deptId := dept.ParentId
	tmpDep := &model.Dept{}
	for true {
		level++
		err = d.orm.Model(&model.Dept{}).Where("dept_id = ?", deptId).First(tmpDep).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, code.MedErrDBFailed
		}
		deptId = tmpDep.ParentId
		resp.List = append(resp.List, &medCalendar.GetTopDeptByUserIdResp_DeptInfo{
			Level: level,
			Dept: &medCalendar.GetTopDeptByUserIdResp_Dept{
				DeptId:   tmpDep.DeptId,
				DeptName: tmpDep.Name,
			},
		})
		parentId := strings.TrimSpace(tmpDep.ParentId)
		if parentId == "1" || parentId == "" {
			break
		}
		tmpDep = &model.Dept{}
		// 防止脏数据
		if level > 6 {
			break
		}
	}
	return resp, nil
}

func (d *Dao) GetDeptData(where map[string]interface{}) (depts []*model.Dept, err error) {
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.Dept{}).Where(whereSql, params...).Find(&depts).Error
	return
}

//
func (d *Dao) FetchOneUser(ctx context.Context, filter map[string]interface{}) (resp *model.User, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	err = d.orm.Model(&model.User{}).Where(whereSql, params...).First(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchOneUser] 查询失败，err：%+v,err")
	}
	return
}

// 查询用户
func (d *Dao) FetchAllUser(ctx context.Context, filter map[string]interface{}) (resp []*model.User, err error) {
	resp = make([]*model.User, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	err = d.orm.Model(&model.User{}).Where(whereSql, params...).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllUser] 查询失败，err：%+v,err")
	}
	return
}

func (d *Dao) UpdateUser(ctx context.Context, data *model.User, filter map[string]interface{}) (err error) {

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	err = d.orm.Model(&model.User{}).Where(whereSql, params...).Updates(data).Error
	if err != nil {
		log.Errorc(ctx, "[dao][[UpdateUser]修改部门，err:%v", err)
	}
	return
}

func (d *Dao) FetchUserDept(ctx context.Context, filter map[string]interface{}) (resp []*model.Dept, err error) {
	resp = make([]*model.Dept, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	err = d.orm.Table("user as u").
		Select("d.*").
		Joins("left join dept as d on d.dept_code = u.dept_code").
		Where(whereSql, params...).
		Order("d.level, d.dept_code asc").
		Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchUserDept] 查询失败，err：%+v,err")
	}
	return
}

// 查询用户
func (d *Dao) GetUserData(filter map[string]interface{}) (resp []*model.User, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.User{}).Where(whereSql, params...).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Error("[dao][GetUserData] 查询失败，err：%+v,err")
	}
	return
}

// 查询用户映射
func (d *Dao) FetchAllUserMap() (resp map[string]*model.User, err error) {
	resp = make(map[string]*model.User, 0)
	var users []*model.User
	err = d.orm.Model(&model.User{}).Find(&users).Error
	err = dbutil.IgnoreNoRecordErr(err)
	if err != nil {
		return
	}

	for _, user := range users {
		resp[user.UserId] = user
	}
	return
}

// 查询用户
func (d *Dao) AddUser(items []*model.User) (err error) {
	err = d.orm.Model(&model.User{}).Create(&items).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Error("[dao][GetUserData] 查询失败，err：%+v,err")
	}
	return
}

// 创建用户
func (d *Dao) CreateUser(ctx context.Context, user *model.User) (err error) {
	err = d.orm.Create(user).Error
	if err != nil {
		log.Errorc(ctx, "[dao][CreateUser] 创建用户，err：%+v,", err)
	}
	return
}

// 获取多个部门用户的数据
func (d *Dao) GetUserDeptDum(deptCode string) (userIds []string, err error) {
	rows, err := d.orm.Raw("select user_id from user group by user_id having count(*) > 1").Rows()
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	userIds = make([]string, 0)
	for rows.Next() {
		var userId string
		if err = rows.Scan(&userId); err != nil {
			return
		}
		userIds = append(userIds, userId)
	}
	if len(userIds) == 0 {
		return
	}
	var users []*model.User
	err = d.orm.Model(&model.User{}).Where("user_id in (?)", userIds).Where("dept_code like '" + deptCode + "%'").Find(&users).Error
	userIds = make([]string, 0)
	userMap := make(map[string]int)
	for _, item := range users {
		if _, ok := userMap[item.UserId]; !ok {
			userIds = append(userIds, item.UserId)
		}
	}
	return
}

// 获取用户的所有一级部门code
func (d *Dao) GetUserTopDeptCodes(ctx context.Context, userId string) []string {
	var rows []model.User
	r := make([]string, 0)

	_ = d.orm.Model(&model.User{}).Where("user_id = ?", userId).Where("state = ?", 1).Find(&rows).Error
	if len(rows) > 0 {
		for _, v := range rows {
			codes := strings.Split(v.DeptCode, "-")
			if len(codes) >= 2 {
				r = append(r, strings.Join(codes[0:2], "-"))
			}
		}
	}

	return r
}

func (d *Dao) GetUserIsLeaderMap(ctx context.Context, params []string) map[string]bool {
	r := make(map[string]bool)
	if len(params) == 0 {
		return r
	}

	userIds := make([]string, 0)

	for _, v := range params {
		r[v] = false
		_, userId, ok := toutils.ParseDeptUidKey(v)
		if ok {
			userIds = append(userIds, fmt.Sprintf("'%s'", userId))
		}
	}
	if len(userIds) == 0 {
		return r
	}

	sql := fmt.Sprintf(`
select concat(user_id, "||" ,dept_code) as uk from user where user_id in (%s) and is_leader = 1 and state = 1
`, strings.Join(userIds, ","))

	type row struct {
		Uk string `gorm:"column:uk"`
	}

	var rows []row

	_ = d.orm.Model(&model.User{}).Raw(sql).Find(&rows).Error

	if len(rows) > 0 {
		for _, v := range rows {
			r[v.Uk] = true
		}
	}

	return r
}
