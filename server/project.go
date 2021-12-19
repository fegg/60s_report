package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	medCalendar "med-common/app/service/med-calendar/api/v1"
	"med-common/app/service/med-calendar/internal/service/config"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/pkg/errors"

	cdb "git.medlinker.com/service/common/db/v2"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// 查询项目
func (d *Dao) GetProjectByName(name string) (project *model.Project, err error) {
	project = &model.Project{}
	db := d.orm.Model(&model.Project{}).Where("project_name = ?", name).
		Where("state > ?", 0)
	err = db.First(project).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func (d *Dao) GetProjectCountForUpdate(projectId int64, name string) (count int64, err error) {
	err = d.orm.Model(&model.Project{}).Where("id != ?", projectId).Where("project_name = ?", name).Where("state != ?", -1).Count(&count).Error
	return
}

func (d *Dao) GetProjectById(ctx context.Context, projectId int64) (project *model.Project, err error) {
	project = new(model.Project)
	err = d.orm.Model(&model.Project{}).Where("id = ?", projectId).Where("state != ?", -1).First(project).Error
	return
}

// 创建项目详情
func (d *Dao) CreateProject(ctx context.Context, userId string, project *model.Project, risks []*model.ProjectRisk, focusUsers []*model.UserFocusProject, checkUserIds []string, checkList []*model.ProjectCheckStd) (err error) {
	updateTask := func(db *gorm.DB) error {
		coin := d.getProjectCoinByLevel(project.ProjectLevel)
		project.Coin = coin

		if err = db.Create(project).Error; err != nil {
			log.Errorc(ctx, "CreateProject err <params:%+v> <err:%v>", project, err)
			err = code.MedErrDBFailed
			return err
		}

		// 扣除父项目coin
		if project.ParentId > 0 {
			efCount := db.Model(&model.Project{}).Where("id = ?", project.ParentId).Where("coin - ? >= 0 ", coin).
				Updates(map[string]interface{}{"coin": gorm.Expr("coin - ?", coin)}).RowsAffected
			if efCount == 0 {
				err = code.MedErrMedCoinNotEnough
				return err
			}
		} else { // 消耗总办的分
			efCount := db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).Where("coin - ? >= 0 ", coin).
				Updates(map[string]interface{}{"coin": gorm.Expr("coin - ?", coin)}).RowsAffected
			if efCount == 0 {
				err = code.MedErrMedAgentCoinNotEnough
				return err
			}
		}

		// 生成结束节点
		node := &model.ProjectNode{
			ProjectId:    project.Id,
			Name:         "项目结束节点",
			FinishTime:   project.FinishTime,
			Importance:   5,
			Status:       2,
			IsAutoCreate: 1,
		}
		if err = db.Create(node).Error; err != nil {
			log.Errorc(ctx, "CreateProject err <params:%+v> <err:%v>", project, err)
			err = code.MedErrDBFailed
			return err
		}

		if len(risks) > 0 {
			for i, _ := range risks {
				risks[i].ProjectId = project.Id
			}
			if err = db.Create(&risks).Error; err != nil {
				log.Errorc(ctx, "CreateProjectRisk err <params:%+v> <err:%v>", risks, err)
				err = code.MedErrDBFailed
				return err
			}
		}

		if len(focusUsers) > 0 {
			for i, _ := range focusUsers {
				focusUsers[i].ProjectId = project.Id
			}
			if err = db.Create(&focusUsers).Error; err != nil {
				log.Errorc(ctx, "CreateUserFocusProject err <params:%+v> <err:%v>", risks, err)
				err = code.MedErrDBFailed
				return err
			}
		}
		if len(checkList) > 0 {
			for i, _ := range checkList {
				checkList[i].ProjectId = project.Id
			}
			if err = db.Create(&checkList).Error; err != nil {
				log.Errorc(ctx, "CreateCheckListProject err <params:%+v> <err:%v>", risks, err)
				err = code.MedErrDBFailed
				return err
			}
			checkUsers := make([]*model.ProjectCheckStdChecker, 0)
			logs := make([]*model.ProjectCheckStdLog, 0)
			for _, item := range checkList {
				bt, _ := json.Marshal(item)
				logs = append(logs, &model.ProjectCheckStdLog{
					CheckId:    item.Id,
					UserId:     userId,
					ProjectId:  item.ProjectId,
					OldContent: "",
					NewContent: string(bt),
					Tag:        1,
					CreatedAt:  time.Now(),
				})
				for _, userId := range checkUserIds {
					checkUsers = append(checkUsers, &model.ProjectCheckStdChecker{
						CheckId:    item.Id,
						UserId:     userId,
						ProjectId:  item.ProjectId,
						CheckState: 0,
						State:      1,
					})
				}
			}
			if err = db.Create(&logs).Error; err != nil {
				log.Errorc(ctx, "CreateCheckListProject err <params:%+v> <err:%v>", risks, err)
				err = code.MedErrDBFailed
				return err
			}
			if len(checkUsers) > 0 {
				if err = db.Create(&checkUsers).Error; err != nil {
					log.Errorc(ctx, "CreateCheckListProject err <params:%+v> <err:%v>", risks, err)
					err = code.MedErrDBFailed
					return err
				}
			}
		}
		return nil
	}
	if err = cdb.DoTrans(d.orm, updateTask); err != nil {
		log.Errorc(ctx, "更新失败 <err:%v>", err)
		//err = code.MedErrDBFailed
		return err
	}
	return
}

// 创建项目详情
func (d *Dao) BatchCreateProject(ctx context.Context, projects []*model.Project) (err error) {
	updateTask := func(db *gorm.DB) error {
		if err = db.Create(projects).Error; err != nil {
			log.Errorc(ctx, "CreateProject err <params:%+v> <err:%v>", projects, err)
			err = code.MedErrDBFailed
			return err
		}
		nodes := make([]*model.ProjectNode, 0)
		for _, item := range projects {
			nodes = append(nodes, &model.ProjectNode{
				ProjectId:    item.Id,
				Name:         "项目结束节点",
				FinishTime:   item.FinishTime,
				Importance:   5,
				Status:       2,
				IsAutoCreate: 1,
			})
		}
		if err = db.Create(nodes).Error; err != nil {
			log.Errorc(ctx, "CreateProject err <params:%+v> <err:%v>", projects, err)
			err = code.MedErrDBFailed
			return err
		}
		return nil
	}
	if err = cdb.DoTrans(d.orm, updateTask); err != nil {
		log.Errorc(ctx, "更新失败 <err:%v>", err)
		err = code.MedErrDBFailed
		return err
	}
	return
}

// 更新项目
func (d *Dao) UpdateProject(ctx context.Context, dingdingUserId string, project *model.Project, req *medCalendar.UpdateProjectReq, ctime time.Time) (newCheck []*model.ProjectCheckStdChecker, err error) {
	newCheck = make([]*model.ProjectCheckStdChecker, 0)
	// 父项目不能是自身
	log.Infoc(ctx, "UpdateProject,project:%+v,req:%+v", project, req)
	if req.ParentId == project.Id {
		err = errors.New("父项目不能设置为自身")
		return
	}
	projectId := project.Id

	finishTime, err := time.Parse("2006-01-02", req.FinishTime)
	if err != nil {
		return
	}

	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return
	}
	officerIds := make([]string, 0)
	if len(req.OfficerIds) > 0 {
		officerIds = req.OfficerIds
	}
	officerIdBt, _ := json.Marshal(officerIds)
	commentUserIds := make([]string, 0)
	if len(req.CommentUserIds) > 0 {
		commentUserIds = req.CommentUserIds
	}
	commentUserIdBt, _ := json.Marshal(commentUserIds)
	checkUserIds := make([]string, 0)
	if len(req.CheckUserIds) > 0 {
		checkUserIds = req.CheckUserIds
	}
	checkUserIdBt, _ := json.Marshal(checkUserIds)

	projectMap := make(map[string]interface{})
	projectMap["parent_id"] = req.ParentId
	projectMap["project_name"] = req.ProjectName
	projectMap["project_rank_type"] = req.ProjectRankType
	projectMap["project_level"] = req.ProjectLevel
	projectMap["manager_id"] = req.ManagerId
	projectMap["manager_name"] = req.ManagerName
	projectMap["project_link"] = req.ProjectLink
	projectMap["cost_dept_code"] = req.CostDeptCode
	projectMap["finish_time"] = finishTime
	projectMap["start_time"] = startTime
	projectMap["auth_control"] = req.AuthControl
	projectMap["officer_ids"] = string(officerIdBt)
	projectMap["evaluate_user_id"] = req.EvaluateUserId
	projectMap["comment_user_ids"] = commentUserIdBt
	projectMap["check_user_ids"] = checkUserIdBt
	projectMap["intro"] = req.Intro
	projectMap["project_subject_id"] = req.ProjectSubjectId
	projectMap["is_water_project"] = req.IsWaterProject
	projectMap["coin"] = d.getProjectCoinByLevel(req.ProjectLevel)
	projectMap["pmo_user_id"] = req.PmoUserId
	projectMap["bi_user_id"] = req.BiUserId
	projectMap["is_cross_dept"] = req.IsCrossDept
	projectMap["cross_dept_codes"] = req.CrossDeptCodes
	var (
		// 需要更新的风险记录
		needAddRisks []model.ProjectRisk
	)
	for _, riskReq := range req.Risks {
		risk := model.ProjectRisk{
			ProjectId:        projectId,
			RiskLevel:        riskReq.RiskLevel,
			RiskDesc:         riskReq.RiskDesc,
			RiskStatus:       riskReq.RiskStatus,
			RelationNodeAble: riskReq.RelationNodeAble,
		}
		needAddRisks = append(needAddRisks, risk)
	}
	updateCheckIds := make([]int64, 0)
	for _, item := range req.CheckList {
		if item.Id > 0 {
			updateCheckIds = append(updateCheckIds, item.Id)
		}
	}
	checkContentMap := make(map[int64]*model.ProjectCheckStd)
	if len(updateCheckIds) > 0 {
		var checkContentList []*model.ProjectCheckStd
		checkContentList, err = d.GetProjectStdData(map[string]interface{}{"project_id": req.ProjectId, "state": 1})
		if err != nil {
			return
		}
		for _, item := range checkContentList {
			checkContentMap[item.Id] = item
		}
		for _, id := range updateCheckIds {
			if _, ok := checkContentMap[id]; !ok {
				return
			}
		}
	}
	var checkUserList []*model.ProjectCheckStdChecker
	checkUserList, err = d.GetProjectCheckerData(req.ProjectId)
	checkUserMap := make(map[int64]map[string]*model.ProjectCheckStdChecker)
	for _, item := range checkUserList {
		if _, ok := checkUserMap[item.CheckId]; !ok {
			checkUserMap[item.CheckId] = make(map[string]*model.ProjectCheckStdChecker)
		}
		checkUserMap[item.CheckId][item.UserId] = item
	}
	updateTask := func(db *gorm.DB) error {
		// 原dbProject 信息
		dbProject, err := d.GetProjectById(ctx, project.Id)
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return fmt.Errorf("获取项目失败 <err:%v>", err)
		}
		// 更新project
		err = db.Model(&model.Project{}).Where("id = ?", projectId).Updates(projectMap).Error
		if err != nil {
			log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("更新项目失败 <err:%v>", err)
		}

		// 比较等级是否变化了
		srcLevel := dbProject.ProjectLevel
		tarLevel := req.ProjectLevel
		if srcLevel != tarLevel {
			if err = dbutil.IgnoreNoRecordErr(err); err != nil {
				return fmt.Errorf("获取父级项目失败 <err:%v>", err)
			}
			// 项目等级对应coin
			srcCoin := d.getProjectCoinByLevel(srcLevel)
			tarCoin := d.getProjectCoinByLevel(tarLevel)
			fmt.Println("-------", srcLevel, tarLevel, srcCoin, tarCoin)
			if srcCoin > tarCoin { // 降级
				// 增加放水项目的curCoin
				if project.ParentId > 0 {
					efCount := db.Model(&model.Project{}).Where("id = ?", project.ParentId).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin + ?", srcCoin-tarCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedCoinNotEnough
					}
				} else {
					efCount := db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin + ?", srcCoin-tarCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedAgentCoinNotEnough
					}
				}
			} else if srcCoin < tarCoin { // 升级
				// 减少放水项目curCoin
				if project.ParentId > 0 {
					efCount := db.Model(&model.Project{}).Where("id = ?", project.ParentId).
						Where("coin - ? >= 0", tarCoin-srcCoin).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin - ?", tarCoin-srcCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedCoinNotEnough
					}
				} else {
					efCount := db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
						Where("coin - ? >= 0", tarCoin-srcCoin).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin - ?", tarCoin-srcCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedAgentCoinNotEnough
					}
				}

			}
		}

		// 删除原来的风险
		err = db.Model(&model.ProjectRisk{}).Where("project_id = ?", req.ProjectId).Delete(&model.ProjectRisk{}).Error
		if err != nil {
			log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("清理当前项目额外风险失败 <err:%v>", err)
		}

		// 新增风险
		if len(needAddRisks) > 0 {
			err := db.Model(&model.ProjectRisk{}).Create(&needAddRisks).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("新增项目风险失败 <err:%v>", err)
			}
		}

		// 更新关注leader
		var leaders []*model.UserFocusProject
		for _, userId := range req.FocusLeaders {
			leaders = append(leaders, &model.UserFocusProject{
				ProjectId: projectId,
				UserId:    userId,
				State:     1,
			})
		}

		// 删除项目关注leader
		err = db.Model(&model.UserFocusProject{}).Where("project_id = ?", req.ProjectId).Delete(&model.UserFocusProject{}).Error
		if err != nil {
			log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
		}

		// 新增项目关注leader
		if len(leaders) > 0 {
			err = db.Create(&leaders).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}

		addContentList := make([]*model.ProjectCheckStd, 0)
		updateContentList := make([]*model.ProjectCheckStd, 0)
		addCheckUsers := make([]*model.ProjectCheckStdChecker, 0)
		dealCheckIds := make([]int64, 0)
		updateCheckIds := make([]int64, 0)
		logs := make([]*model.ProjectCheckStdLog, 0)
		cMap := make(map[int64]int)
		for _, item := range req.CheckList {
			if item.Id == 0 {
				content := &model.ProjectCheckStd{
					ProjectId:    req.ProjectId,
					CheckTime:    item.CheckTime,
					Percent:      item.Percent,
					CheckContent: item.CheckContent,
					State:        1,
					CheckState:   10,
					CreatedAt:    ctime,
					CreateTime:   ctime.Unix(),
				}
				addContentList = append(addContentList, content)
			} else {
				cMap[item.Id] = 1
				var isChange bool
				if checkContentMap[item.Id].CheckTime != item.CheckTime || checkContentMap[item.Id].Percent != item.Percent || checkContentMap[item.Id].CheckContent != item.CheckContent {
					if checkContentMap[item.Id].VerifyState != 0 {
						return fmt.Errorf("审核中不能修改内容 <err:%v>", err)
					}
					obt, _ := json.Marshal(checkContentMap[item.Id])
					checkContentMap[item.Id].CheckTime = item.CheckTime
					checkContentMap[item.Id].Percent = item.Percent
					checkContentMap[item.Id].CheckContent = item.CheckContent
					checkContentMap[item.Id].State = 1
					checkContentMap[item.Id].CheckState = 10
					checkContentMap[item.Id].CreatedAt = ctime
					checkContentMap[item.Id].CreateTime = ctime.Unix()
					updateContentList = append(updateContentList, checkContentMap[item.Id])
					nbt, _ := json.Marshal(checkContentMap[item.Id])
					logs = append(logs, &model.ProjectCheckStdLog{
						UserId:     dingdingUserId,
						ProjectId:  req.ProjectId,
						OldContent: string(obt),
						NewContent: string(nbt),
						Tag:        2,
						CheckId:    item.Id,
						CreatedAt:  time.Now(),
						CreateTime: ctime.Unix(),
					})
					isChange = true
				} else {
					// 已审核通过
					if checkContentMap[item.Id].CheckState == 20 {
						continue
					}
				}
				for _, userId := range req.CheckUserIds {
					if _, ok := checkUserMap[item.Id]; ok {
						if _, ok := checkUserMap[item.Id][userId]; !ok {
							addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
								CheckId:   item.Id,
								UserId:    userId,
								ProjectId: projectId,
								State:     1,
							})
						} else {
							if isChange {
								updateCheckIds = append(updateCheckIds, checkUserMap[item.Id][userId].Id)
							}
						}
					} else {
						addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
							CheckId:   item.Id,
							UserId:    userId,
							ProjectId: projectId,
							State:     1,
						})
					}
				}
				for userId, citem := range checkUserMap[item.Id] {
					var isMatch bool
					for _, cuserId := range req.CheckUserIds {
						if userId == cuserId {
							isMatch = true
							break
						}
					}
					if !isMatch {
						dealCheckIds = append(dealCheckIds, citem.Id)
					}
				}
			}
		}
		dIds := []int64{}
		for id := range checkContentMap {
			if _, ok := cMap[id]; !ok {
				dIds = append(dIds, id)
			}
		}
		if len(dIds) > 0 {
			err = db.Model(&model.ProjectCheckStd{}).Where("id in(?)", dIds).Update("is_delete", 1).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(addContentList) > 0 {
			err = db.Create(addContentList).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
			for _, item := range addContentList {
				for _, userId := range req.CheckUserIds {
					addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
						CheckId:   item.Id,
						UserId:    userId,
						ProjectId: projectId,
						State:     1,
					})
				}
				bt, _ := json.Marshal(item)
				logs = append(logs, &model.ProjectCheckStdLog{
					UserId:     dingdingUserId,
					CheckId:    item.Id,
					ProjectId:  req.ProjectId,
					OldContent: "",
					NewContent: string(bt),
					Tag:        1,
					CreatedAt:  time.Now(),
					CreateTime: ctime.Unix(),
				})
			}

		}
		if len(updateContentList) > 0 {
			for _, item := range updateContentList {
				err = db.Model(&model.ProjectCheckStd{}).Where("id = ?", item.Id).Updates(&item).Error
				if err != nil {
					log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
					return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
				}
			}
		}
		if len(addCheckUsers) > 0 {
			newCheck = addCheckUsers
			err = db.Create(addCheckUsers).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(updateCheckIds) > 0 {
			err = db.Model(&model.ProjectCheckStdChecker{}).Where("id in(?)", updateCheckIds).Update("check_state", 0).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		log.Info("aaaaa %v", dealCheckIds)
		if len(dealCheckIds) > 0 {
			err = db.Model(&model.ProjectCheckStdChecker{}).Where("id in(?)", dealCheckIds).Update("state", -1).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(logs) > 0 {
			err = db.Create(logs).Error
			if err != nil {
				log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		return nil
	}
	if err = cdb.DoTrans(d.orm, updateTask); err != nil {
		log.Errorc(ctx, "UpdateProject <err:%v> <req:%+v>", err, req)
		//err = code.MedErrDBFailed
		return
	}
	return
}

func (d *Dao) UpdateProjectFinishStatus(projectId int64, updates map[string]interface{}) (err error) {
	return d.orm.Model(&model.Project{}).Where("id = ?", projectId).Updates(updates).Error
}

func (d *Dao) GetChildrenProject(ctx context.Context, parentId int64) (projects []*model.Project, err error) {
	err = d.orm.Where("parent_id = ?", parentId).
		Where("state = ?", 1).
		Where("project_type = ?", 0).
		Find(&projects).Error
	err = dbutil.IgnoreNoRecordErr(err)
	return
}

func (d *Dao) GetProjectByIds(ctx context.Context, ids []int64) (projectMap map[int64]*model.Project, err error) {
	projectMap = make(map[int64]*model.Project, 0)
	var projects []*model.Project
	err = d.orm.Where("id in (?)", ids).Find(&projects).Error
	for _, p := range projects {
		projectMap[p.Id] = p
	}
	return
}

func (d *Dao) FetchAllProjectMap(ctx context.Context) (projectMap map[int64]*model.Project, err error) {
	projectMap = make(map[int64]*model.Project, 0)
	var projects []*model.Project
	err = d.orm.Find(&projects).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao|project] FetchAllProjectMap select all fail err(%+v)", err)
		return
	}
	for _, p := range projects {
		projectMap[p.Id] = p
	}
	return
}

func (d *Dao) FetchProjectMap(ctx context.Context, where map[string]interface{}) (projectMap map[int64]*model.Project, err error) {
	projectMap = make(map[int64]*model.Project, 0)
	var projects []*model.Project
	db := d.orm
	if len(where) > 0 {
		db = db.Where(where)
	}
	err = db.Find(&projects).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao|project] FetchAllProjectMap select all fail err(%+v)", err)
		return
	}
	for _, p := range projects {
		projectMap[p.Id] = p
	}
	return
}

func (d *Dao) GetProjectMapByIds(ctx context.Context, ids []int64) (projectMap map[int64]*model.Project, err error) {
	projectMap = make(map[int64]*model.Project, 0)
	var projects []*model.Project
	err = d.orm.Where("id in (?)", ids).Find(&projects).Error
	for _, p := range projects {
		projectMap[p.Id] = p
	}
	return
}

// 查询项目日报
func (d *Dao) GetProjectCost(ctx context.Context, projectIds []int64, startDate string, endDate string) (total int64, err error) {
	var cost []int64
	// 查询子项目耗费
	err = d.orm.Model(&model.ProjectTask{}).
		Where("time_scope >= ?", startDate).
		Where("time_scope <= ?", endDate).
		Where("state > ?", 0).
		Where("project_id in (?)", projectIds).
		Pluck("COALESCE(sum(cost),0) as cost", &cost).Error
	if err != gorm.ErrRecordNotFound {
		err = nil
	}
	if len(cost) > 0 {
		total = cost[0]
	}
	return
}

func (d *Dao) GetProjectTask(ctx context.Context, where map[string]interface{}) (reports []*model.ProjectTask, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.ProjectTask{}).Where(whereSql, params...).Find(&reports).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func (d *Dao) GetProjectDeptStatics(ctx context.Context, ids []int64, startTime string, endTime string) (ret map[string]*model.ProjectDeptStaticsItem, err error) {
	db := d.orm.Model(&model.ProjectTask{}).
		Where("state > ?", 0).
		Where("cate in (?)", []int64{1, 2}).
		Where("project_id in (?)", ids)

	if startTime != "" {
		db = db.Where("time_scope >= ?", startTime)
	}

	if endTime != "" {
		db = db.Where("time_scope <= ?", endTime)
	}
	var items []*model.ProjectTask
	err = db.Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	ret = make(map[string]*model.ProjectDeptStaticsItem)
	for _, item := range items {
		deptCode := item.UserDeptCode
		if len(deptCode) > 7 {
			deptCode = deptCode[:7]
		}
		if _, ok := ret[deptCode]; !ok {
			ret[deptCode] = &model.ProjectDeptStaticsItem{
				DeptCode: deptCode,
				Cost:     0,
			}
		}
		ret[deptCode].Cost += item.Cost
	}

	return
}

func (d *Dao) GetProjectUserStatics(ctx context.Context, projectIds []int64, startTime string, endTime string) (ret []*model.ProjectTask, err error) {
	db := d.orm.Model(&model.ProjectTask{}).
		Where("state > ?", 0).
		Where("project_id in (?)", projectIds)

	if startTime != "" {
		db = db.Where("time_scope >= ?", startTime)
	}

	if endTime != "" {
		db = db.Where("time_scope <= ?", endTime)
	}
	err = db.Find(&ret).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (d *Dao) GetDeptUserStatics(ctx context.Context, deptCode string, startTime string, endTime string) (ret int64, err error) {
	var data []*model.ProjectDeptStaticsItem
	db := d.orm.Model(&model.ProjectTask{}).
		Where("state > ?", 0).
		Where("user_dept_code like ?", deptCode+"%")

	if startTime != "" {
		db = db.Where("time_scope >= ?", startTime)
	}

	if endTime != "" {
		db = db.Where("time_scope <= ?", endTime)
	}

	err = db.Scan(&data).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	for _, item := range data {
		ret += item.Cost
	}

	return
}

func (d *Dao) GetDeptCostStatics(ctx context.Context, startTime string, endTime string) (m map[string]int64, err error) {
	db := d.orm.Model(&model.ProjectTask{}).Select("user_dept_code,sum(cost) as cost").Where("state = 1").Group("user_dept_code")
	if startTime != "" {
		db = db.Where("time_scope >= ?", startTime)
	}
	if endTime != "" {
		db = db.Where("time_scope <= ?", endTime)
	}
	rows, err := db.Rows()
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	m = make(map[string]int64)
	defer rows.Close()
	for rows.Next() {
		var deptCode string
		var cost int64
		if err = rows.Scan(&deptCode, &cost); err != nil {
			return
		}
		if len(deptCode) > 7 {
			deptCode = deptCode[0:7]
		}
		m[deptCode] += cost
	}
	return
}

// 部门下所有人的填报时间
func (d *Dao) GetDeptAllUserStaticsCostTotal(ctx context.Context, deptCode string, startTime string, endTime string) (ret int64, err error) {
	var users []*model.User
	// 获取部门下所有的人
	db := d.orm.Model(&model.User{}).
		Select("user_id").
		Where("state > ?", 0).
		Where("dept_code like ?", deptCode+"%")
	if err = db.Find(&users).Error; err != nil {
		log.Errorc(ctx, "[dao|calendar|project] GetDeptAllUserStaticsCostTotal select user err params(%s) error(%+v)", deptCode, err)
		return
	}
	userIds := make([]string, 0)
	for _, v := range users {
		userIds = append(userIds, v.UserId)
	}

	// 获取填报了的人数
	hasUserIds := make([]string, 0)
	db2 := d.orm.Model(&model.UserDailyReport{}).
		Select("user_id").
		Where("state > ?", 0).
		Where("user_id in (?)", userIds).
		Group("user_id")

	if startTime != "" {
		db2 = db2.Where("time_scope >= ?", startTime)
	}

	if endTime != "" {
		db2 = db2.Where("time_scope <= ?", endTime)
	}

	err = db2.Pluck("user_id", &hasUserIds).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	ret = int64(len(hasUserIds) * 100)

	return
}

// 更新项目标签
func (d *Dao) PostProjectTag(ctx context.Context, projectTag *model.ProjectTag) (err error) {
	var dbProjectTag *model.ProjectTag
	if err = d.orm.Model(&model.ProjectTag{}).Where("name = ?", projectTag.Name).Where("state = ?", 1).
		First(&dbProjectTag).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return code.MedErrDBFailed
		}
	}
	// 名字已经存在
	if dbProjectTag != nil && dbProjectTag.Id > 0 {
		if projectTag.Id == 0 { // 新增
			return code.MedErrProjectTagRepeat
		} else { // 修改
			if dbProjectTag.Id != projectTag.Id {
				return code.MedErrProjectTagRepeat
			}
		}
	}
	if projectTag.Id > 0 {
		// 更新
		if err = d.orm.Model(&model.ProjectTag{}).Where("id = ?", projectTag.Id).Where("state = ?", 1).
			Updates(&projectTag).Error; err != nil {
			return code.MedErrDBFailed
		}
		return
	}
	if err = d.orm.Model(&model.ProjectTag{}).Create(&projectTag).Error; err != nil {
		return code.MedErrDBFailed
	}
	return
}

// 删除项目标签
func (d *Dao) DeleteProjectTag(ctx context.Context, tagId int64) (err error) {
	// 更新
	if err = d.orm.Model(&model.ProjectTag{}).Where("id = ?", tagId).
		Updates(map[string]interface{}{"state": -1}).Error; err != nil {
		return code.MedErrDBFailed
	}
	return
}

// 检索项目标签
func (d *Dao) GetProjectTagListPaginate(ctx context.Context, where map[string]interface{}, order []string, page, pageSize int64) (res []*model.ProjectTag, total int64, err error) {
	res = make([]*model.ProjectTag, 0)
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	db := d.orm.Model(&model.ProjectTag{}).Where(whereSql, params...)

	orders := strings.Join(order, ",")
	if orders != "" {
		db = db.Order(orders)
	}
	countDb := db
	countDb.Count(&total)
	err = db.Offset((int(page) - 1) * int(pageSize)).Limit(int(pageSize)).Find(&res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao|project] GetProjectTagList fail err(%+v)", err)
		return
	}
	return
}

// 检索项目标签
func (d *Dao) GetProjectList(ctx context.Context, where map[string]interface{}, order []string) (res []*model.Project, err error) {
	res = make([]*model.Project, 0)
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	db := d.orm.Model(&model.Project{}).Where(whereSql, params...)

	orders := strings.Join(order, ",")
	if orders != "" {
		db = db.Order(orders)
	}
	err = db.Find(&res).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao|project] GetProjectList fail err(%+v)", err)
		return
	}
	return
}

// 获取项目focus map
func (d *Dao) FetchProjectFocusUserIdMap(ctx context.Context, where map[string]interface{}) (res map[int64][]string, err error) {
	res = make(map[int64][]string)

	projectFocus := make([]*model.UserFocusProject, 0)
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	db := d.orm.Model(&model.UserFocusProject{}).Where(whereSql, params...)
	err = db.Order("updated_at desc").Find(&projectFocus).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao|project] GetProjectList fail err(%+v)", err)
		return
	}
	for _, v := range projectFocus {
		if _, ok := res[v.ProjectId]; !ok {
			res[v.ProjectId] = make([]string, 0)
		}
		res[v.ProjectId] = append(res[v.ProjectId], v.UserId)
	}
	return
}

// 放水
func (d *Dao) ProjectWater(ctx context.Context, userId string, projectId int64, coin int64) (err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][ProjectWater] 启动事务失败，err：%+v,", err)
		return errors.Wrap(err, "启动事务失败")
	}
	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][ProjectWater] 项目放水，err：%+v,", err)
			return
		}
	}()

	var dbProject *model.Project
	if err = d.orm.Model(&model.Project{}).Where("id = ?", projectId).Where("state = ?", 1).First(&dbProject).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			log.Errorc(ctx, "[dao][ProjectWater] 获取项目数据错误，err：%+v,", err)
		}
	}
	if dbProject.Id == 0 {
		return code.MedErrProjectNotExist
	}
	if dbProject.Coin+coin < 0 {
		return code.MedErrMedCoinNotEnough
	}

	efCount := new_db.Model(&model.Project{}).Where("id = ?", projectId).Where("coin + ? >= 0", coin).
		Updates(map[string]interface{}{"coin": gorm.Expr("coin + ?", coin)}).RowsAffected
	if efCount == 0 {
		return code.MedErrMedCoinNotEnough
	}

	efCount = new_db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
		Where("coin + ? >= 0", -1*coin).Updates(map[string]interface{}{
		"coin": gorm.Expr("coin + ?", -1*coin),
	}).RowsAffected
	if efCount == 0 {
		return code.MedErrMedAgentCoinNotEnough
	}

	if err = new_db.Model(&model.ProjectWaterLog{}).Create(&model.ProjectWaterLog{
		ProjectId: projectId,
		UserId:    userId,
		Coin:      coin,
		CreatedAt: time.Now(),
	}).Error; err != nil {
		return
	}
	if err = new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][ProjectWater] 提交事务失败，err：%+v,", err)
		return err
		//return errors.Wrap(err, "提交事务失败")
	}

	return
}

// 放水信息
func (d *Dao) ProjectWaterInfo(ctx context.Context) int64 {

	var coins []int64
	d.orm.Model(&model.Project{}).Where("parent_id = ?", 0).Select("sum(coin) coin").Pluck("coin", &coins)
	if len(coins) > 0 {
		return coins[0]
	}
	return 0
}

// 不同项目的花费coin <levelNum:levelCoin>
func (d *Dao) ProjectLevelCoin() map[int64]int64 {
	appCfg := config.GetAppConfig()
	if err := paladin.Get("application.txt").UnmarshalTOML(&appCfg); err != nil {
		log.Info("init server, load application.text failed: %+v\n", err)
		panic(err)
	}
	// 1-C,2-B,3-A,4-S -1-D
	levelCoinMap := make(map[int64]int64)
	if len(appCfg.BusinessCfg.LevelNum) > 0 {
		for k, v := range appCfg.BusinessCfg.LevelNum {
			levelCoinMap[v] = appCfg.BusinessCfg.LevelCoin[k]
		}
	}

	return levelCoinMap
}

// 不同项目对应的等级 <levelName, levelNum>
func (d *Dao) ProjectLevelNameLevelNum() map[string]int64 {
	appCfg := config.GetAppConfig()
	if err := paladin.Get("application.txt").UnmarshalTOML(&appCfg); err != nil {
		log.Info("init server, load application.text failed: %+v\n", err)
		panic(err)
	}
	// 1-C,2-B,3-A,4-S -1-D
	levelNumMap := make(map[string]int64)
	if len(appCfg.BusinessCfg.LevelName) > 0 {
		for k, v := range appCfg.BusinessCfg.LevelName {
			levelNumMap[v] = appCfg.BusinessCfg.LevelNum[k]
		}
	}

	return levelNumMap
}
// <levelName, levelCoin>
func (d *Dao) ProjectLevelNameLevelCoin() map[string]int64 {
	appCfg := config.GetAppConfig()
	if err := paladin.Get("application.txt").UnmarshalTOML(&appCfg); err != nil {
		log.Info("init server, load application.text failed: %+v\n", err)
		panic(err)
	}
	// 1-C,2-B,3-A,4-S -1-D
	levelCoinMap := make(map[string]int64)
	if len(appCfg.BusinessCfg.LevelName) > 0 {
		for k, v := range appCfg.BusinessCfg.LevelName {
			levelCoinMap[v] = appCfg.BusinessCfg.LevelCoin[k]
		}
	}

	return levelCoinMap
}

// 获取等级对应的coin
func (d *Dao) getProjectCoinByLevel(level int64) int64 {
	levelCoin := d.ProjectLevelCoin()
	if coin, ok := levelCoin[level]; ok {
		return coin
	}
	return 0
}

// 获取一级项目下的总数
func (d *Dao) GetChildProjectCoinCostMap() (ret map[int64]int64) {
	ret = make(map[int64]int64)

	var projectData []*model.Project
	var err error
	if err = d.orm.Model(&model.Project{}).Where("state = ?", 1).Find(&projectData).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return
		}
	}

	for _, v := range projectData {
		if _, ok := ret[v.ParentId]; !ok {
			ret[v.ParentId] = 0
		}
		ret[v.ParentId] = ret[v.ParentId] + v.Coin
	}

	return
}

// 事项转移到项目
func (d *Dao) TransferProject(ctx context.Context, dingdingUserId string, project *model.Project, req *medCalendar.TransferProjectReq, ctime time.Time) (newCheck []*model.ProjectCheckStdChecker, err error) {
	newCheck = make([]*model.ProjectCheckStdChecker, 0)
	// 父项目不能是自身
	log.Infoc(ctx, "TransferProject,project:%+v,req:%+v", project, req)
	if req.ParentId == project.Id {
		err = errors.New("父项目不能设置为自身")
		return
	}
	if req.ProjectLevel > 3 { // 部门级事项不能是S
		err = code.MedErrDeptProjectCanNotS
		return
	}

	projectId := project.Id

	finishTime, err := time.Parse("2006-01-02", req.FinishTime)
	if err != nil {
		return
	}

	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return
	}
	officerIds := make([]string, 0)
	if len(req.OfficerIds) > 0 {
		officerIds = req.OfficerIds
	}
	officerIdBt, _ := json.Marshal(officerIds)
	commentUserIds := make([]string, 0)
	if len(req.CommentUserIds) > 0 {
		commentUserIds = req.CommentUserIds
	}
	commentUserIdBt, _ := json.Marshal(commentUserIds)
	checkUserIds := make([]string, 0)
	if len(req.CheckUserIds) > 0 {
		checkUserIds = req.CheckUserIds
	}
	checkUserIdBt, _ := json.Marshal(checkUserIds)

	projectMap := make(map[string]interface{})
	projectMap["parent_id"] = req.ParentId
	projectMap["project_name"] = req.ProjectName
	projectMap["project_rank_type"] = req.ProjectRankType
	projectMap["project_level"] = req.ProjectLevel
	projectMap["manager_id"] = req.ManagerId
	projectMap["manager_name"] = req.ManagerName
	projectMap["project_link"] = req.ProjectLink
	projectMap["cost_dept_code"] = req.CostDeptCode
	projectMap["finish_time"] = finishTime
	projectMap["start_time"] = startTime
	projectMap["auth_control"] = req.AuthControl
	projectMap["officer_ids"] = string(officerIdBt)
	projectMap["evaluate_user_id"] = req.EvaluateUserId
	projectMap["comment_user_ids"] = commentUserIdBt
	projectMap["check_user_ids"] = checkUserIdBt
	projectMap["intro"] = req.Intro
	projectMap["project_subject_id"] = req.ProjectSubjectId
	projectMap["is_water_project"] = req.IsWaterProject
	projectMap["coin"] = d.getProjectCoinByLevel(req.ProjectLevel)
	projectMap["pmo_user_id"] = req.PmoUserId
	projectMap["bi_user_id"] = req.BiUserId
	projectMap["project_type"] = 0 // 默认为项目
	project.ParentId = req.ParentId
	project.ProjectRankType = req.ProjectRankType
	project.ProjectType = 0 // 默认为项目
	project.ProjectLevel = req.ProjectLevel
	var (
		// 需要更新的风险记录
		needAddRisks []model.ProjectRisk
	)
	for _, riskReq := range req.Risks {
		risk := model.ProjectRisk{
			ProjectId:        projectId,
			RiskLevel:        riskReq.RiskLevel,
			RiskDesc:         riskReq.RiskDesc,
			RiskStatus:       riskReq.RiskStatus,
			RelationNodeAble: riskReq.RelationNodeAble,
		}
		needAddRisks = append(needAddRisks, risk)
	}
	updateCheckIds := make([]int64, 0)
	for _, item := range req.CheckList {
		if item.Id > 0 {
			updateCheckIds = append(updateCheckIds, item.Id)
		}
	}
	checkContentMap := make(map[int64]*model.ProjectCheckStd)
	if len(updateCheckIds) > 0 {
		var checkContentList []*model.ProjectCheckStd
		checkContentList, err = d.GetProjectStdData(map[string]interface{}{"project_id": req.ProjectId, "state": 1})
		if err != nil {
			return
		}
		for _, item := range checkContentList {
			checkContentMap[item.Id] = item
		}
		for _, id := range updateCheckIds {
			if _, ok := checkContentMap[id]; !ok {
				return
			}
		}
	}
	var checkUserList []*model.ProjectCheckStdChecker
	checkUserList, err = d.GetProjectCheckerData(req.ProjectId)
	checkUserMap := make(map[int64]map[string]*model.ProjectCheckStdChecker)
	for _, item := range checkUserList {
		if _, ok := checkUserMap[item.CheckId]; !ok {
			checkUserMap[item.CheckId] = make(map[string]*model.ProjectCheckStdChecker)
		}
		checkUserMap[item.CheckId][item.UserId] = item
	}
	updateTask := func(db *gorm.DB) error {
		// 原dbProject 信息
		dbProject, err := d.GetProjectById(ctx, project.Id)
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return fmt.Errorf("获取项目失败 <err:%v>", err)
		}
		// 更新project
		err = db.Model(&model.Project{}).Where("id = ?", projectId).Updates(projectMap).Error
		if err != nil {
			log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("更新项目失败 <err:%v>", err)
		}

		// 比较等级是否变化了
		srcLevel := dbProject.ProjectLevel
		tarLevel := req.ProjectLevel
		if srcLevel != tarLevel {
			if err = dbutil.IgnoreNoRecordErr(err); err != nil {
				return fmt.Errorf("获取父级项目失败 <err:%v>", err)
			}
			// 项目等级对应coin
			var srcCoin int64 = 0 // 原来为事项 没有等级
			tarCoin := d.getProjectCoinByLevel(tarLevel)
			fmt.Println("-------", srcLevel, tarLevel, srcCoin, tarCoin)
			if srcCoin > tarCoin { // 降级
				// 增加放水项目的curCoin
				if project.ParentId > 0 {
					efCount := db.Model(&model.Project{}).Where("id = ?", project.ParentId).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin + ?", srcCoin-tarCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedCoinNotEnough
					}
				} else {
					efCount := db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin + ?", srcCoin-tarCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedAgentCoinNotEnough
					}
				}
			} else if srcCoin < tarCoin { // 升级
				// 减少放水项目curCoin
				if project.ParentId > 0 {
					efCount := db.Model(&model.Project{}).Where("id = ?", project.ParentId).
						Where("coin - ? >= 0", tarCoin-srcCoin).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin - ?", tarCoin-srcCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedCoinNotEnough
					}
				} else {
					efCount := db.Model(&model.OfficeWater{}).Where("year = ?", time.Now().Year()).
						Where("coin - ? >= 0", tarCoin-srcCoin).
						Updates(map[string]interface{}{
							"coin": gorm.Expr("coin - ?", tarCoin-srcCoin),
						}).RowsAffected
					if efCount == 0 {
						return code.MedErrMedAgentCoinNotEnough
					}
				}
			}
		}
		cate := d.getProjectCate(project)
		if cate == 0 {
			return code.MedErrProjectTypeErr
		}

		// 更新填报类型
		db.Model(&model.ProjectTask{}).Where("project_id = ?", project.Id).
			Updates(map[string]interface{}{
				"cate": cate,
			})

		// 删除原来的风险
		err = db.Model(&model.ProjectRisk{}).Where("project_id = ?", req.ProjectId).Delete(&model.ProjectRisk{}).Error
		if err != nil {
			log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("清理当前项目额外风险失败 <err:%v>", err)
		}

		// 新增风险
		if len(needAddRisks) > 0 {
			err := db.Model(&model.ProjectRisk{}).Create(&needAddRisks).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("新增项目风险失败 <err:%v>", err)
			}
		}

		// 更新关注leader
		var leaders []*model.UserFocusProject
		for _, userId := range req.FocusLeaders {
			leaders = append(leaders, &model.UserFocusProject{
				ProjectId: projectId,
				UserId:    userId,
				State:     1,
			})
		}

		// 删除项目关注leader
		err = db.Model(&model.UserFocusProject{}).Where("project_id = ?", req.ProjectId).Delete(&model.UserFocusProject{}).Error
		if err != nil {
			log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
			return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
		}

		// 新增项目关注leader
		if len(leaders) > 0 {
			err = db.Create(&leaders).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}

		addContentList := make([]*model.ProjectCheckStd, 0)
		updateContentList := make([]*model.ProjectCheckStd, 0)
		addCheckUsers := make([]*model.ProjectCheckStdChecker, 0)
		dealCheckIds := make([]int64, 0)
		updateCheckIds := make([]int64, 0)
		logs := make([]*model.ProjectCheckStdLog, 0)
		cMap := make(map[int64]int)
		for _, item := range req.CheckList {
			if item.Id == 0 {
				content := &model.ProjectCheckStd{
					ProjectId:    req.ProjectId,
					CheckTime:    item.CheckTime,
					Percent:      item.Percent,
					CheckContent: item.CheckContent,
					State:        1,
					CheckState:   10,
					CreatedAt:    ctime,
					CreateTime:   ctime.Unix(),
				}
				addContentList = append(addContentList, content)
			} else {
				cMap[item.Id] = 1
				var isChange bool
				if checkContentMap[item.Id].CheckTime != item.CheckTime || checkContentMap[item.Id].Percent != item.Percent || checkContentMap[item.Id].CheckContent != item.CheckContent {
					if checkContentMap[item.Id].VerifyState != 0 {
						return fmt.Errorf("审核中不能修改内容 <err:%v>", err)
					}
					obt, _ := json.Marshal(checkContentMap[item.Id])
					checkContentMap[item.Id].CheckTime = item.CheckTime
					checkContentMap[item.Id].Percent = item.Percent
					checkContentMap[item.Id].CheckContent = item.CheckContent
					checkContentMap[item.Id].State = 1
					checkContentMap[item.Id].CheckState = 10
					checkContentMap[item.Id].CreatedAt = ctime
					checkContentMap[item.Id].CreateTime = ctime.Unix()
					updateContentList = append(updateContentList, checkContentMap[item.Id])
					nbt, _ := json.Marshal(checkContentMap[item.Id])
					logs = append(logs, &model.ProjectCheckStdLog{
						UserId:     dingdingUserId,
						ProjectId:  req.ProjectId,
						OldContent: string(obt),
						NewContent: string(nbt),
						Tag:        2,
						CheckId:    item.Id,
						CreatedAt:  time.Now(),
						CreateTime: ctime.Unix(),
					})
					isChange = true
				} else {
					// 已审核通过
					if checkContentMap[item.Id].CheckState == 20 {
						continue
					}
				}
				for _, userId := range req.CheckUserIds {
					if _, ok := checkUserMap[item.Id]; ok {
						if _, ok := checkUserMap[item.Id][userId]; !ok {
							addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
								CheckId:   item.Id,
								UserId:    userId,
								ProjectId: projectId,
								State:     1,
							})
						} else {
							if isChange {
								updateCheckIds = append(updateCheckIds, checkUserMap[item.Id][userId].Id)
							}
						}
					} else {
						addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
							CheckId:   item.Id,
							UserId:    userId,
							ProjectId: projectId,
							State:     1,
						})
					}
				}
				for userId, citem := range checkUserMap[item.Id] {
					var isMatch bool
					for _, cuserId := range req.CheckUserIds {
						if userId == cuserId {
							isMatch = true
							break
						}
					}
					if !isMatch {
						dealCheckIds = append(dealCheckIds, citem.Id)
					}
				}
			}
		}
		dIds := []int64{}
		for id := range checkContentMap {
			if _, ok := cMap[id]; !ok {
				dIds = append(dIds, id)
			}
		}
		if len(dIds) > 0 {
			err = db.Model(&model.ProjectCheckStd{}).Where("id in(?)", dIds).Update("is_delete", 1).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(addContentList) > 0 {
			err = db.Create(addContentList).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
			for _, item := range addContentList {
				for _, userId := range req.CheckUserIds {
					addCheckUsers = append(addCheckUsers, &model.ProjectCheckStdChecker{
						CheckId:   item.Id,
						UserId:    userId,
						ProjectId: projectId,
						State:     1,
					})
				}
				bt, _ := json.Marshal(item)
				logs = append(logs, &model.ProjectCheckStdLog{
					UserId:     dingdingUserId,
					CheckId:    item.Id,
					ProjectId:  req.ProjectId,
					OldContent: "",
					NewContent: string(bt),
					Tag:        1,
					CreatedAt:  time.Now(),
					CreateTime: ctime.Unix(),
				})
			}

		}
		if len(updateContentList) > 0 {
			for _, item := range updateContentList {
				err = db.Model(&model.ProjectCheckStd{}).Where("id = ?", item.Id).Updates(&item).Error
				if err != nil {
					log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
					return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
				}
			}
		}
		if len(addCheckUsers) > 0 {
			newCheck = addCheckUsers
			err = db.Create(addCheckUsers).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(updateCheckIds) > 0 {
			err = db.Model(&model.ProjectCheckStdChecker{}).Where("id in(?)", updateCheckIds).Update("check_state", 0).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		log.Info("aaaaa %v", dealCheckIds)
		if len(dealCheckIds) > 0 {
			err = db.Model(&model.ProjectCheckStdChecker{}).Where("id in(?)", dealCheckIds).Update("state", -1).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		if len(logs) > 0 {
			err = db.Create(logs).Error
			if err != nil {
				log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
				return fmt.Errorf("删除项目关注leader失败 <err:%v>", err)
			}
		}
		return nil
	}
	if err = cdb.DoTrans(d.orm, updateTask); err != nil {
		log.Errorc(ctx, "TransferProject <err:%v> <req:%+v>", err, req)
		//err = code.MedErrDBFailed
		return
	}
	return
}

func (d *Dao) getProjectCate(project *model.Project) (ret int64) {
	ret = 0
	switch project.ProjectRankType {
	case 1:
		ret = 1
		break
	case 2:
		ret = 2
		break
	case 0:
		switch project.ProjectType {
		case 1:
			ret = 3
		case 2:
			ret = 4
		}
	}
	return
}

func (d *Dao) GetProjectFullById(ctx context.Context, id int64) *model.ProjectFull {
	if id == 0 {
		return nil
	}

	var row model.ProjectFull

	err := d.orm.Model(&model.ProjectFull{}).Where("id = ?", id).First(&row).Error
	if err != nil {
		log.Errorc(ctx, "GetProjectFullById id: %d, err: %+v", id, err)
		return nil
	}

	return &row
}

func (d *Dao) HasValidChildProject(ctx context.Context, parentId int64) bool {
	var r int64
	err := d.orm.Model(&model.ProjectFull{}).Where("parent_id = ?", parentId).Where("state = ?", 1).Count(&r).Error
	if err != nil {
		log.Errorc(ctx, "HasValidChildProject id: %d, err: %+v", parentId, err)
		return false
	}

	return r > 0
}

func (d *Dao) InsertProjectDeleteLog(ctx context.Context, dat *model.ProjectDeleteLog) (insId int64) {

	err := d.orm.Model(&model.ProjectDeleteLog{}).Create(dat).Error

	if err != nil {
		log.Errorc(ctx, "InsertProjectDeleteLog dat: %+v, err: %+v", dat, err)
		return
	}
	insId = dat.Id

	return
}

func (d *Dao) SetProjectDeleted(ctx context.Context, projectId int64) bool {
	rows := d.orm.Model(&model.ProjectFull{}).Where("id = ?", projectId).Where("state = ?", 1).Update("state", -1).RowsAffected

	log.Errorc(ctx, "SetProjectDeleted id: %+v, rows: %+v", projectId, rows)

	return rows > 0
}

func (d *Dao) IncreaseProjectCoin(ctx context.Context, projectId, coin int64) bool {
	_ = d.orm.Model(&model.ProjectFull{}).Where("id = ?", projectId).Update("coin", coin).Error

	return true
}
