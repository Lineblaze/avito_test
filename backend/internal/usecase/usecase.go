package usecase

import (
	"fmt"
	openapi "github.com/Lineblaze/avito_gen"
	repository "zadanie-6105/backend/internal"
	"zadanie-6105/backend/internal/domain"
)

//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s UseCase -p internal -y "Controller describes methods, implemented by the usecase package."
type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Employees

func (u *UseCase) GetEmployeeByID(employeeID int64) (*domain.Employee, error) {
	employee, err := u.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("getting employee by ID: %w", err)
	}
	return employee, nil
}

func (u *UseCase) GetEmployeeByUserName(username string) (*domain.Employee, error) {
	employee, err := u.repo.GetEmployeeByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("getting employee by username: %w", err)
	}
	return employee, nil
}

func (u *UseCase) CreateEmployee(req *domain.CreateEmployeeRequest) (*domain.Employee, error) {
	employeeInfo := domain.Employee{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	createdEmployee, err := u.repo.CreateEmployee(&employeeInfo)
	if err != nil {
		return nil, fmt.Errorf("creating employee: %w", err)
	}

	return createdEmployee, nil
}

// Organizations

func (u *UseCase) GetOrganizationByID(organizationID int64) (*domain.Organization, error) {
	organization, err := u.repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, fmt.Errorf("getting organization by ID: %w", err)
	}
	return organization, nil
}

func (u *UseCase) CreateOrganization(req *domain.CreateOrganizationRequest) (*domain.Organization, error) {
	organizationInfo := domain.Organization{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	createdOrganization, err := u.repo.CreateOrganization(&organizationInfo)
	if err != nil {
		return nil, fmt.Errorf("creating organization: %w", err)
	}

	return createdOrganization, nil
}

func (u *UseCase) AssignEmployeeToOrganization(req *domain.AssignEmployeeToOrganizationRequest) (*domain.OrganizationResponsible, error) {
	orgRespInfo := domain.OrganizationResponsible{
		OrganizationID: req.OrganizationID,
		UserID:         req.UserID,
	}
	assign, err := u.repo.AssignEmployeeToOrganization(&orgRespInfo)
	if err != nil {
		return nil, fmt.Errorf("assigning employee to organization: %w", err)
	}
	return assign, nil
}

func (u *UseCase) CheckUserOrganizationResponsibility(organizationId string) (bool, error) {
	return u.repo.IsUserResponsibleForOrganization(organizationId)
}

// Tenders

func (u *UseCase) GetTenders() ([]openapi.Tender, error) {
	tenders, err := u.repo.GetTenders()
	if err != nil {
		return nil, fmt.Errorf("getting tenders: %w", err)
	}
	return tenders, nil
}

func (u *UseCase) GetUserTenders(userName string) ([]*openapi.Tender, error) {
	tenders, err := u.repo.GetUserTenders(userName)
	if err != nil {
		return nil, fmt.Errorf("getting tenders: %w", err)
	}
	return tenders, nil
}

func (u *UseCase) GetTenderStatus(tenderID string) (string, error) {
	status, err := u.repo.GetTenderStatus(tenderID)
	if err != nil {
		return "", fmt.Errorf("getting tender status: %w", err)
	}
	return status, nil
}

func (u *UseCase) CreateTender(req *openapi.CreateTenderRequest) (*openapi.Tender, error) {
	tender := openapi.Tender{
		Name:           req.Name,
		Description:    req.Description,
		ServiceType:    req.ServiceType,
		Status:         "Created",
		OrganizationId: req.OrganizationId,
		Version:        1,
	}

	createdTender, err := u.repo.CreateTender(&tender)
	if err != nil {
		return nil, fmt.Errorf("creating tender: %w", err)
	}

	return createdTender, nil
}

func (u *UseCase) EditTender(tenderID string, req *openapi.EditTenderRequest) (*openapi.Tender, error) {
	existingTender, err := u.repo.GetTenderByID(tenderID)
	if err != nil {
		return nil, fmt.Errorf("getting existing tender: %w", err)
	}

	if req.Name != nil {
		existingTender.Name = *req.Name
	}
	if req.Description != nil {
		existingTender.Description = *req.Description
	}
	if req.ServiceType != nil {
		existingTender.ServiceType = *req.ServiceType
	}

	existingTender.Version += 1

	updatedTender, err := u.repo.EditTender(existingTender)
	if err != nil {
		return nil, fmt.Errorf("updating tender: %w", err)
	}
	return updatedTender, nil
}

func (u *UseCase) UpdateTenderStatus(tenderID string, status string) error {
	err := u.repo.UpdateTenderStatus(tenderID, status)
	if err != nil {
		return fmt.Errorf("updating tender status: %w", err)
	}
	return nil
}

func (u *UseCase) RollbackTender(tenderID string, version string) (*openapi.Tender, error) {
	tenderAtVersion, err := u.repo.GetTenderByVersion(tenderID, version)
	if err != nil {
		return nil, fmt.Errorf("retrieving tender version: %w", err)
	}

	updatedTender, err := u.repo.EditTender(tenderAtVersion)
	if err != nil {
		return nil, fmt.Errorf("updating tender to rolled-back version: %w", err)
	}

	return updatedTender, nil
}

// Bids
