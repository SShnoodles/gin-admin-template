package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var (
	ErrMenuPathExists    = errors.New("menu path exists")
	ErrMenuNotFound      = errors.New("menu not found")
	ErrInvalidMenuParent = errors.New("invalid menu parent")
)

type MenuTree struct {
	domain.Menu
	Children []*MenuTree `json:"children"`
}

func FindMenuByPath(path string) (domain.Menu, error) {
	var menu domain.Menu
	err := config.DB.First(&menu, "path = ?", path).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func FindMenusByPid(pid int64) ([]domain.Menu, error) {
	var menu []domain.Menu
	err := config.DB.Find(&menu, "pid = ?", pid).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func DeleteMenusByPid(pid int64) error {
	return config.DB.Delete(&domain.Menu{}, "pid = ?", pid).Error
}

func FindResourceIdsByMenuId(id int64) ([]string, error) {
	var mr []domain.MenuResourceRelation
	err := config.DB.Find(&mr, "menu_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	if len(mr) == 0 {
		return ids, nil
	}
	for _, m := range mr {
		ids = append(ids, strconv.FormatInt(m.ResourceId, 10))
	}
	return ids, nil
}

func FindMenuTree() ([]*MenuTree, error) {
	var menus []domain.Menu
	err := FindAll(&menus)
	if err != nil {
		return nil, err
	}
	var menuTree []*MenuTree
	for _, menu := range menus {
		var mt MenuTree
		copier.Copy(&mt, menu)
		menuTree = append(menuTree, &mt)
	}
	return buildMenuTree(menuTree, 0, make(map[int64]bool)), nil
}

func buildMenuTree(menuTree []*MenuTree, pid int64, visited map[int64]bool) []*MenuTree {
	var children []*MenuTree
	for _, menu := range menuTree {
		if menu.Pid == pid {
			if visited[menu.Id] {
				config.Log.Warnf("skip circular menu reference, id=%d pid=%d", menu.Id, menu.Pid)
				continue
			}
			visited[menu.Id] = true
			menu.Children = buildMenuTree(menuTree, menu.Id, visited)
			delete(visited, menu.Id)
			children = append(children, menu)
		}
	}
	sort.SliceStable(children, func(i, j int) bool {
		if children[i].Sort == children[j].Sort {
			return children[i].Id < children[j].Id
		}
		return children[i].Sort < children[j].Sort
	})
	return children
}

func CreateMenu(menu domain.Menu, resourceIds []string) (int64, error) {
	if err := validateMenu(menu); err != nil {
		return 0, err
	}
	_, err := FindMenuByPath(menu.Path)
	if err == nil {
		return 0, ErrMenuPathExists
	}

	menuId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		menu.Id = menuId
		if err = tx.Create(&menu).Error; err != nil {
			return err
		}
		return replaceMenuResources(tx, menuId, resourceIds)
	})
	if err != nil {
		return 0, err
	}
	return menuId, nil
}

func UpdateMenu(menuId int64, input domain.Menu, resourceIds []string) error {
	if err := validateMenu(input); err != nil {
		return err
	}
	input.Id = menuId
	var menu domain.Menu
	err := FindById(&menu, menuId)
	if err != nil {
		return ErrMenuNotFound
	}
	if input.Path != menu.Path {
		_, err = FindMenuByPath(input.Path)
		if err == nil {
			return ErrMenuPathExists
		}
	}
	if input.Pid == menuId {
		return ErrInvalidMenuParent
	}
	if input.Pid != 0 {
		var menus []domain.Menu
		err = FindAll(&menus)
		if err != nil {
			return err
		}
		if isMenuDescendant(menus, menuId, input.Pid, make(map[int64]bool)) {
			return ErrInvalidMenuParent
		}
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		copier.Copy(&menu, input)
		if err = tx.Save(&menu).Error; err != nil {
			return err
		}
		return replaceMenuResources(tx, menuId, resourceIds)
	})
}

func DeleteMenu(id int64) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Menu{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("pid = ?", id).Delete(&domain.Menu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", id).Delete(&domain.MenuResourceRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func replaceMenuResources(tx *gorm.DB, menuId int64, resourceIds []string) error {
	if err := tx.Where("menu_id = ?", menuId).Delete(&domain.MenuResourceRelation{}).Error; err != nil {
		return err
	}
	if len(resourceIds) == 0 {
		return nil
	}
	ids, err := ParsePositiveIds(resourceIds)
	if err != nil {
		return err
	}
	var mrr []domain.MenuResourceRelation
	for _, resourceId := range ids {
		mrr = append(mrr, domain.MenuResourceRelation{
			Id:         config.IdGenerate(),
			ResourceId: resourceId,
			MenuId:     menuId,
		})
	}
	return tx.Create(&mrr).Error
}

func validateMenu(menu domain.Menu) error {
	if menu.Pid < 0 || strings.TrimSpace(menu.Name) == "" || strings.TrimSpace(menu.Path) == "" || strings.TrimSpace(menu.Icon) == "" || menu.Sort < 0 {
		return ErrInvalidParam
	}
	return nil
}

func isMenuDescendant(menus []domain.Menu, rootId int64, targetId int64, visited map[int64]bool) bool {
	for _, menu := range menus {
		if menu.Pid != rootId {
			continue
		}
		if visited[menu.Id] {
			continue
		}
		if menu.Id == targetId {
			return true
		}
		visited[menu.Id] = true
		if isMenuDescendant(menus, menu.Id, targetId, visited) {
			return true
		}
	}
	return false
}
