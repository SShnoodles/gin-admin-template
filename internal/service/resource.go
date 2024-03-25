package service

import (
	"encoding/json"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"github.com/go-openapi/spec"
	"os"
)

func SaveResourceFromSwagger(docPath string) {
	file, err := os.ReadFile(docPath)
	if err != nil {
		config.Log.Error(err.Error())
		return
	}
	var sw spec.Swagger
	if err := json.Unmarshal(file, &sw); err != nil {
		config.Log.Error(err.Error())
		return
	}
	for path, pathItem := range sw.Paths.Paths {

		for method, operation := range map[string]*spec.Operation{
			"get":    pathItem.Get,
			"post":   pathItem.Post,
			"put":    pathItem.Put,
			"delete": pathItem.Delete,
		} {
			if operation != nil {
				resource := domain.Resource{
					Name:   operation.Summary,
					Method: method,
					Path:   path,
				}
				oldResource, _ := FindResourceByMethodAndPath(method, path)
				if oldResource == (domain.Resource{}) {
					resource.Id = config.IdGenerate()
					Insert(&resource)
				} else {
					resource.Id = oldResource.Id
					Update(&resource)
				}
			}
		}
	}
}

func FindResourceByMethodAndPath(method, path string) (domain.Resource, error) {
	var resource domain.Resource
	err := config.DB.First(&resource, "method = ? and path = ?", method, path).Error
	if err != nil {
		return resource, err
	}
	return resource, nil
}

func FindResourcesByUserId(id int64) ([]domain.Resource, error) {
	var resources []domain.Resource

	var user domain.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return resources, err
	}

	var roles []domain.UserRoleRelation
	err = config.DB.Find(&roles, "user_id = ?", id).Error
	if err != nil {
		return resources, err
	}
	var roleIds []int64
	for _, role := range roles {
		roleIds = append(roleIds, role.RoleId)
	}

	var menus []domain.RoleMenuRelation
	err = config.DB.Find(&menus, "role_id IN ?", roleIds).Error
	if err != nil {
		return resources, err
	}
	var menuIds []int64
	for _, menu := range menus {
		menuIds = append(menuIds, menu.MenuId)
	}

	var mrs []domain.MenuResourceRelation
	err = config.DB.Find(&mrs, "menu_id IN ?", menuIds).Error
	if err != nil {
		return resources, err
	}
	var resourceIds []int64
	for _, resource := range mrs {
		resourceIds = append(resourceIds, resource.ResourceId)
	}

	err = config.DB.Find(&resources, "id IN ?", resourceIds).Error
	if err != nil {
		return resources, err
	}
	return resources, nil
}
