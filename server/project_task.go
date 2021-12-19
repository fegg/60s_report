package calendar

import (
	"context"
	"med-common/app/service/med-calendar/internal/pkg/time_util"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"sync"
	"time"

	"github.com/go-kratos/kratos/pkg/sync/errgroup"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (d *Dao) FetchProjectTaskPaginateGroupByUid(ctx context.Context, filter map[string]interface{}, groups []string, page, limit int) (resp []string, err error) {
	resp = make([]string, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	for _, group := range groups {
		db = db.Group(group)
	}
	err = db.Order("project_task.updated_at desc").Limit(limit).Offset((page-1)*limit).Pluck("user_id", &resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchProjectTaskPaginate] 查询失败，err：%+v,err", err)
	}
	return
}

func (d *Dao) GetProjectTaskPaginateGroupByUid(ctx context.Context, filter map[string]interface{}, groups []string, page, limit int) (resp []string, err error) {
	resp = make([]string, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Joins("left join user on project_task.user_id = user.user_id").Where(whereSql, params...)
	for _, group := range groups {
		db = db.Group(group)
	}
	err = db.Order("project_task.updated_at desc").Limit(limit).Offset((page-1)*limit).Pluck("project_task.user_id", &resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetProjectTaskPaginateGroupByUid] 查询失败，err：%+v,err", err)
	}
	return
}

func (d *Dao) FetchProjectTaskCountGroupByUid(ctx context.Context, filter map[string]interface{}, groups []string) (total int64, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Select("count(DISTINCT user_id)").Where(whereSql, params...)
	for _, group := range groups {
		db = db.Group(group)
	}
	if err = db.Count(&total).Error; err != nil {
		log.Errorc(ctx, "[dao][FetchProjectTaskCount] 查询count失败，err：%+v,err", err)
		return
	}

	return
}

func (d *Dao) GetProjectTaskCountGroupByUid(ctx context.Context, filter map[string]interface{}, groups []string) (total int64, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Select("count(DISTINCT project_task.user_id)").Joins("left join user on project_task.user_id = user.user_id").Where(whereSql, params...)
	for _, group := range groups {
		db = db.Group(group)
	}
	if err = db.Count(&total).Error; err != nil {
		log.Errorc(ctx, "[dao][GetProjectTaskCountGroupByUid] 查询count失败，err：%+v,err", err)
		return
	}

	return
}

// 查询项目
func (d *Dao) FetchAllProjectTask(ctx context.Context, filter map[string]interface{}) (project []*model.ProjectTask, err error) {
	project = make([]*model.ProjectTask, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	err = db.Order("updated_at desc").Find(&project).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllProjectTask] 查询失败，err：%+v,", err)
	}
	return
}

// 查询项目
func (d *Dao) GetAllProjectTask(ctx context.Context, filter map[string]interface{}) (project []*model.ProjectTask, err error) {
	project = make([]*model.ProjectTask, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Select("project_task.*").Joins("left join user on project_task.user_id = user.user_id").Where(whereSql, params...)
	err = db.Order("project_task.updated_at desc").Find(&project).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllProjectTask] 查询失败，err：%+v,", err)
	}
	return
}

// 查询项目
func (d *Dao) FetchOneProjectTask(ctx context.Context, filter map[string]interface{}, order string) (project *model.ProjectTask, err error) {
	project = &model.ProjectTask{}
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	if order != "" {
		db = db.Order(order)
	}
	err = db.Order("id desc").First(project).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchOneProjectTask] 查询失败，err：%+v,", err)
	}
	return
}

// 删除
func (d *Dao) DeleteProjectTask(ctx context.Context, filter map[string]interface{}) (err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Where(whereSql, params...)
	err = db.Delete(&model.ProjectTask{}).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][DeleteProjectTask] 删除每日填报，err：%+v,", err)
	}
	return
}

// 修改
func (d *Dao) UpdateProjectTask(ctx context.Context, data *model.ProjectTask, filter map[string]interface{}) (err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	err = db.Updates(data).Error
	if err != nil {
		log.Errorc(ctx, "[dao][[UpdateProjectTask]修改每日填报失败，err:%v", err)
	}
	return
}

// 创建
func (d *Dao) CreateProjectTask(ctx context.Context, project *model.ProjectTask) (err error) {
	db := d.orm.Model(&model.ProjectTask{})
	err = db.Create(project).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][DeleteProjectTask] 删除每日填报，err：%+v,", err)
	}
	return
}

// 批量创建
func (d *Dao) BatchCreateProjectTask(ctx context.Context, projects []*model.ProjectTask) (err error) {
	db := d.orm.Model(&model.ProjectTask{})
	err = db.Create(&projects).Error
	if err != nil {
		log.Errorc(ctx, "[dao][BatchCreateProjectTask] 创建每日填报，err：%+v,", err)
	}
	return
}

// 批量修改或者创建每日填报、以及统计
func (d *Dao) TransCreateOrUpdateProjectTask(ctx context.Context, isAutoFill bool, insert []*model.ProjectTask, userDailyData []*model.UserDailyReport, userId, timeScope, checkUserId string) (err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][BatchCreateOrUpdateProjectTask] 启动事务失败，err：%+v,", err)
		return errors.Wrap(err, "启动事务失败")
	}

	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][BatchCreateOrUpdateProjectTask] 修改数据失败，err：%+v,", err)
			return
		}
	}()

	// 查询出校准前的数据
	checkBeforeJson := "{}"
	//删除日报
	if err = new_db.Model(&model.ProjectTask{}).Where("user_id = ?", userId).
		Where("time_scope = ?", timeScope).Update("state", -1).Error; err != nil {
		return
	}

	// 查询删除的所有日报
	delProjectTaskIds := make([]int64, 0)
	if err = new_db.Model(&model.ProjectTask{}).Where("user_id = ?", userId).
		Where("time_scope = ?", timeScope).Pluck("id", &delProjectTaskIds).Error; err != nil {
		return
	}

	if err = new_db.Model(&model.UserDailyReport{}).Where("user_id = ?", userId).
		Where("time_scope = ?", timeScope).Update("state", -1).Error; err != nil {
		return
	}
	// 插入日报
	if len(insert) > 0 {
		if err = new_db.Model(&model.ProjectTask{}).Create(&insert).Error; err != nil {
			return
		}
	}

	// 插入统计
	times, _ := time.ParseInLocation("2006-01-02", timeScope, time.Local)
	if len(userDailyData) > 0 {
		if err = new_db.Model(&model.UserDailyReport{}).Create(&userDailyData).Error; err != nil {
			return
		}
	} else {
		// 获取用户数据
		var userInfo *model.User
		if err = new_db.Model(&model.User{}).Where("user_id = ?", userId).First(&userInfo).Error; err != nil {
			return
		}
		// 插入百分百损耗
		if err = new_db.Model(&model.UserDailyReport{}).Create(&model.UserDailyReport{
			UserId:             userId,
			UserDeptCode:       userInfo.DeptCode,
			CompanyProjectCost: 0,
			DeptProjectCost:    0,
			MatterCost:         0,
			ManageCost:         0,
			Loss:               100,
			IsLeave:            0,
			IsAutoFill:         0,
			State:              1,
			TimeScope:          times,
			CreatedAt:          time.Now(),
		}).Error; err != nil {
			return
		}
	}

	// 查询出校准后的数据
	checkAfterJson := "{}"
	// 记录校准日志
	//times, _ := time.ParseInLocation("2006-01-02", timeScope, time.Local)
	log.Infoc(ctx, "[dao|calendar|project_task] TransCreateOrUpdateProjectTask ProjectTaskCheckLog info (%s, %s, %+v)", checkBeforeJson, checkAfterJson, times)

	weekDay := time_util.GetCurrentMondayDate(0).Format("2006-01-02")
	var autoFill []*model.UserAutoFill
	err = new_db.Model(&model.UserAutoFill{}).Where("user_id = ? and week_day = ?", userId, weekDay).Find(&autoFill).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if isAutoFill {
		if len(autoFill) > 0 {
			err = new_db.Model(&model.UserAutoFill{}).Where("id = ?", autoFill[0].Id).Update("is_auto_fill", 1).Error
		} else {
			err = new_db.Model(&model.UserAutoFill{}).Create(&model.UserAutoFill{UserId: userId, WeekDay: weekDay, IsAutoFill: 1, State: 1}).Error
		}
	} else {
		if len(autoFill) > 0 && autoFill[0].IsAutoFill > 0 {
			err = new_db.Model(&model.UserAutoFill{}).Where("id = ?", autoFill[0].Id).Update("is_auto_fill", 0).Error
		}
	}
	if err != nil {
		return
	}

	if err := new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][TransCreateOrUpdateProjectTask] 提交事务失败，err：%+v,", err)
		return errors.Wrap(err, "提交事务失败")
	}
	return
}

// 查询项目
type StatsDailyReportEntity struct {
	TotalCost int `gorm:"column:total_cost" column:"total_cost" json:"total_cost"`
	model.ProjectTask
}

func (d *Dao) GetUserDailyReport(where map[string]interface{}) (reports []*model.UserDailyReport, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	err = d.orm.Model(&model.UserDailyReport{}).Where(whereSql, params...).Find(&reports).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func (d *Dao) GetUserDailyReportList(where map[string]interface{}, start, limit int) (reports []*model.UserDailyReport, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.UserDailyReport{}).Where(whereSql, params...).Order("id").Offset(start).Limit(limit).Find(&reports).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func (d *Dao) AutoCreateProjectTask(ctx context.Context, userIds []string, userDailyData []*model.UserDailyReport, insert []*model.ProjectTask) (err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][AutoCreateProjectTask] 启动事务失败，err：%+v,", err)
		return errors.Wrap(err, "启动事务失败")
	}

	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][AutoCreateProjectTask] 修改数据失败，err：%+v,", err)
			return
		}
	}()
	if len(userIds) > 0 {
		err = new_db.Model(&model.UserDailyReport{}).Where("user_id in (?) and time_scope = ?", userIds, time.Now().Format("2006-01-02")).Update("state", -1).Error
		if err != nil {
			return
		}
		err = new_db.Model(&model.ProjectTask{}).Where("user_id in (?) and time_scope = ?", userIds, time.Now().Format("2006-01-02")).Update("state", -1).Error
		if err != nil {
			return
		}
	}
	// 插入统计
	if err = new_db.Model(&model.UserDailyReport{}).Create(&userDailyData).Error; err != nil {
		return
	}
	if len(insert) > 0 {
		// 插入日报
		if err = new_db.Model(&model.ProjectTask{}).Create(&insert).Error; err != nil {
			return
		}
	}
	if err := new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][BatchCreateOrUpdateProjectTask] 提交事务失败，err：%+v,", err)
		return errors.Wrap(err, "提交事务失败")
	}
	return
}

type UserProjectTaskInfo struct {
	UserId    string `gorm:"user_id"`
	ProjectId int64  `gorm:"project_id"`
	Cost      int64  `gorm:"cost"`
}

// 获取用户某个时间段内的填报数据
func (d *Dao) GetUserProjectTaskInfo(ctx context.Context, userId string, startTime, endTime string) (res []*UserProjectTaskInfo, err error) {
	filter := map[string]interface{}{
		"user_id": userId,
		"state":   1,
	}
	if startTime != "" && endTime != "" {
		filter["time_scope"] = []string{dbutil.WhereConditionStringBetween, startTime, endTime}
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	if err = db.Select("user_id, project_id, sum(cost) cost").Group("user_id, project_id").Find(&res).Error; err != nil {
		return
	}
	return
}

func (d *Dao) GetLastUserProjectTask(ctx context.Context, where map[string]interface{}) (m map[string][]*model.ProjectTask, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	rows, err := d.orm.Model(&model.ProjectTask{}).Select("max(time_scope) time_scope, user_id").Where(whereSql, params...).Group("user_id").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	m = make(map[string][]*model.ProjectTask)
	ids := make(map[string][][]string, 0)

	userTimeMap := make(map[string][]string)
	for rows.Next() {
		var timeScope time.Time
		var userId string
		if err = rows.Scan(&timeScope, &userId); err != nil {
			return
		}
		day := timeScope.Format("2006-01-02")
		userTimeMap[day] = append(userTimeMap[day], userId)
		if len(userTimeMap[day]) == 200 {
			ids[day] = append(ids[day], userTimeMap[day])
			userTimeMap[day] = make([]string, 0)
		}
	}
	for day, userIds := range userTimeMap {
		if len(userIds) > 0 {
			ids[day] = append(ids[day], userTimeMap[day])
		}
	}
	if len(ids) == 0 {
		return
	}
	var metux sync.Mutex
	eg := errgroup.WithContext(ctx)
	for day, userIds := range ids {
		for _, userId := range userIds {
			func(timeScope string, userIds []string) {
				eg.Go(func(ctx context.Context) error {
					items, err := d.GetProjectTask(ctx, map[string]interface{}{"state": 1, "time_scope": timeScope, "user_id": userIds})
					if err != nil {
						return err
					}
					metux.Lock()
					for _, item := range items {
						if item.ProjectId < 0 {
							continue
						}
						m[item.UserId] = append(m[item.UserId], item)
					}
					metux.Unlock()
					return nil
				})
			}(day, userId)
		}
	}
	err = eg.Wait()
	return
}

// 获取用户指定时间段内是否填写
func (d *Dao) GetUserFillLogDuration(ctx context.Context, userId, startTime, endTime string) map[string]bool {
	r := make(map[string]bool)

	type row struct {
		TimeScope time.Time `gorm:"column:time_scope"`
	}

	var rows []row

	err := d.orm.Model(&model.ProjectTask{}).
		Select("time_scope").
		Where("user_id = ?", userId).
		Where("state = ?", 1).
		Where("time_scope >= ?", startTime).
		Where("time_scope <= ?", endTime).
		Group("time_scope").
		Order("time_scope desc").
		Find(&rows).Error

	if err != nil {
		log.Errorc(ctx, "GetUserFillLogDuration err: %+v", err)
		return r
	}

	for _, v := range rows {
		r[v.TimeScope.Format("2006-01-02")] = true
	}

	return r
}

// 某个项目是否被填写过
func (d *Dao) HasProjectBeenFilled(ctx context.Context, projectId int64) bool {
	var cnt int64

	_ = d.orm.Model(&model.ProjectTask{}).
		Where("state = ?", 1).
		Where("project_id = ?", projectId).
		Count(&cnt).
		Error

	return cnt > 0
}


// 转移工作量
func (d *Dao) TransferWorkload(ctx context.Context, projectId, afterProjectId int64) (err error){
	err = d.orm.Model(&model.ProjectTask{}).Where("project_id = ?", projectId).
		Updates(map[string]interface{}{
			"project_id": afterProjectId,
			"before_project_id": projectId,
	}).Error
	if err != nil {
		log.Errorc(ctx, "[dao|project_task] TransferWorkload fail err(%+v)", err)
		return
	}
	return
}

// 获取部门项目的所有耗时
func (d *Dao) ProjectDeptCostMap(ctx context.Context, filter map[string]interface{}) (resp map[string][][]int64, err error) {
	resp = make(map[string][][]int64)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...)
	rows, err := db.Select("project_id, cost_dept_code, sum(cost) cost").Group("project_id, cost_dept_code").Rows()
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][ProjectCostMap] 查询失败，err：%+v,err", err)
	}
	for rows.Next() {
		type item struct {
			ProjectId int64 `json:"project_id"`
			CostDeptCode string `json:"cost_dept_code"`
			Cost int64 `json:"cost"`
		}
		var i item
		db.ScanRows(rows, &i)
		if _, ok := resp[i.CostDeptCode]; !ok {
			resp[i.CostDeptCode] = make([][]int64, 0)
		}
		resp[i.CostDeptCode] = append(resp[i.CostDeptCode], []int64{i.ProjectId, i.Cost})
	}
	return
}