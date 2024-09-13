package usecase

import (
	"fmt"
	openapi "github.com/Lineblaze/avito_gen"
	repository "zadanie-6105/backend/internal"
)

//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s UseCase -p internal -y "Controller describes methods, implemented by the usecase package."
type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Employees

func (u *UseCase) GetEmployeeByID(employeeID int64) (*openapi.Employee, error) {
	employee, err := u.repo.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, fmt.Errorf("getting employee by ID: %w", err)
	}
	return employee, nil
}

func (u *UseCase) GetEmployeeByUserName(username string) (*openapi.Employee, error) {
	employee, err := u.repo.GetEmployeeByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("getting employee by username: %w", err)
	}
	return employee, nil
}

func (u *UseCase) CreateEmployee(req *openapi.CreateEmployeeRequest) (*openapi.Employee, error) {
	employeeInfo := openapi.Employee{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	}

	createdEmployee, err := u.repo.CreateEmployee(&employeeInfo)
	if err != nil {
		return nil, fmt.Errorf("creating employee: %w", err)
	}

	return createdEmployee, nil
}

// Organizations

func (u *UseCase) GetOrganizationByID(organizationID int64) (*openapi.Organization, error) {
	organization, err := u.repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, fmt.Errorf("getting organization by ID: %w", err)
	}
	return organization, nil
}

func (u *UseCase) CreateOrganization(req *openapi.CreateOrganizationRequest) (*openapi.Organization, error) {
	organizationInfo := openapi.Organization{
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

func (u *UseCase) AssignEmployeeToOrganization(req *openapi.AssignEmployeeToOrganizationRequest) (*openapi.OrganizationResponsible, error) {
	orgRespInfo := openapi.OrganizationResponsible{
		OrganizationId: req.OrganizationId,
		UserId:         req.UserId,
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

func (u *UseCase) CheckUserOrganizationResponsibilityByUsername(username string) (bool, error) {
	return u.repo.IsUserResponsibleForOrganizationByUsername(username)
}

// Tenders

func (u *UseCase) GetTenders() ([]openapi.Tender, error) {
	tenders, err := u.repo.GetTenders()
	if err != nil {
		return nil, fmt.Errorf("getting tenders: %w", err)
	}
	return tenders, nil
}

func (u *UseCase) GetTenderByID(tenderID string) (*openapi.Tender, error) {
	tender, err := u.repo.GetTenderByID(tenderID)
	if err != nil {
		return nil, fmt.Errorf("getting bid by ID: %w", err)
	}
	return tender, nil
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

func (u *UseCase) CanUserAccessTender(tenderID string) (bool, error) {
	tender, err := u.repo.GetTenderByID(tenderID)
	if err != nil {
		return false, fmt.Errorf("failed to get tender: %w", err)
	}

	switch tender.Status {
	case "Created", "Closed":
		isResponsible, err := u.repo.IsUserResponsibleForOrganization(tender.OrganizationId)
		if err != nil {
			return false, fmt.Errorf("failed to check user responsibility: %w", err)
		}
		return isResponsible, nil

	case "Open", "Published":
		return true, nil

	default:
		return false, nil
	}
}

// Bids

func (u *UseCase) GetBidByID(bidID string) (*openapi.Bid, error) {
	bid, err := u.repo.GetBidByID(bidID)
	if err != nil {
		return nil, fmt.Errorf("getting bid by ID: %w", err)
	}
	return bid, nil
}

func (u *UseCase) GetUserBids(userName string) ([]*openapi.Bid, error) {
	bids, err := u.repo.GetUserBids(userName)
	if err != nil {
		return nil, fmt.Errorf("getting tenders: %w", err)
	}
	return bids, nil
}

func (u *UseCase) GetBidsByTenderID(tenderID string) ([]*openapi.Bid, error) {
	bids, err := u.repo.GetBidsByTenderID(tenderID)
	if err != nil {
		return nil, fmt.Errorf("getting bids by tender ID: %w", err)
	}
	return bids, nil
}

func (u *UseCase) GetBidStatus(bidID string) (string, error) {
	status, err := u.repo.GetBidStatus(bidID)
	if err != nil {
		return "", fmt.Errorf("getting bid status: %w", err)
	}
	return status, nil
}

func (u *UseCase) CreateBid(req *openapi.CreateBidRequest) (*openapi.Bid, error) {
	exists, err := u.repo.BidExistsByTenderID(req.TenderId)
	if err != nil {
		return nil, fmt.Errorf("error checking existing bid for tender %s: %w", req.TenderId, err)
	}

	if exists {
		return nil, fmt.Errorf("a bid for this tender %s already exists", req.TenderId)
	}

	bid := openapi.Bid{
		Name:        req.Name,
		Description: req.Description,
		Status:      "Created",
		TenderId:    req.TenderId,
		AuthorId:    req.OrganizationId,
		AuthorType:  "Organization",
		Version:     1,
	}

	createdBid, err := u.repo.CreateBid(&bid)
	if err != nil {
		return nil, fmt.Errorf("creating bid: %w", err)
	}

	return createdBid, nil
}

func (u *UseCase) EditBid(bidID string, req *openapi.EditBidRequest) (*openapi.Bid, error) {
	existingBid, err := u.repo.GetBidByID(bidID)
	if err != nil {
		return nil, fmt.Errorf("getting existing bid: %w", err)
	}

	if req.Name != nil {
		existingBid.Name = *req.Name
	}
	if req.Description != nil {
		existingBid.Description = *req.Description
	}

	existingBid.Version += 1

	updatedBid, err := u.repo.EditBid(existingBid)
	if err != nil {
		return nil, fmt.Errorf("updating bid: %w", err)
	}
	return updatedBid, nil
}

func (u *UseCase) UpdateBidStatus(bidID string, status string) error {
	err := u.repo.UpdateBidStatus(bidID, status)
	if err != nil {
		return fmt.Errorf("updating bid status: %w", err)
	}
	return nil
}

func (u *UseCase) RollbackBid(bidID string, version string) (*openapi.Bid, error) {
	bidAtVersion, err := u.repo.GetBidByVersion(bidID, version)
	if err != nil {
		return nil, fmt.Errorf("retrieving bid version: %w", err)
	}

	updatedBid, err := u.repo.EditBid(bidAtVersion)
	if err != nil {
		return nil, fmt.Errorf("updating bid to rolled-back version: %w", err)
	}

	return updatedBid, nil
}

func (u *UseCase) SubmitBidDecision(bidId string, decision string, username string) error {
	err := u.repo.UpdateBidDecision(bidId, decision, username)
	if err != nil {
		return fmt.Errorf("error updating bid decision for bid %s: %w", bidId, err)
	}

	decisions, err := u.repo.GetBidDecisions(bidId)
	if err != nil {
		return fmt.Errorf("error getting bid decisions for bid %s: %w", bidId, err)
	}

	for _, d := range decisions {
		if d == "Rejected" {
			err = u.repo.RejectBid(bidId)
			if err != nil {
				return fmt.Errorf("error rejecting bid %s: %w", bidId, err)
			}
			return nil
		}
	}

	responsibles, err := u.repo.GetResponsibleUsersForOrganization()
	if err != nil {
		return fmt.Errorf("error getting responsible users for organization: %w", err)
	}

	quorum := len(responsibles)

	approvedCount := 0
	for _, d := range decisions {
		if d == "Approved" {
			approvedCount++
		}
	}

	if approvedCount >= quorum && quorum >= 3 {
		err = u.repo.CloseTenderByBid(bidId)
		if err != nil {
			return fmt.Errorf("error closing tender after bid approval for bid %s: %w", bidId, err)
		}
	}

	return nil
}

func (u *UseCase) IsTenderClosed(bidId string) (bool, error) {
	tenderStatus, err := u.repo.GetTenderStatusByBid(bidId)
	if err != nil {
		return false, fmt.Errorf("error getting tender status for bid %s: %w", bidId, err)
	}
	return tenderStatus == "Closed", nil
}

func (u *UseCase) SubmitBidFeedback(bidId string, feedback string, username string) error {
	err := u.repo.UpdateBidFeedback(bidId, feedback, username)
	if err != nil {
		return fmt.Errorf("error updating feedback for bid %s: %w", bidId, err)
	}
	return nil
}

func (u *UseCase) GetBidReviewsByTenderId(tenderId string) ([]openapi.BidReview, error) {
	reviews, err := u.repo.GetBidReviewsByTenderId(tenderId)
	if err != nil {
		return nil, fmt.Errorf("error fetching reviews for tender %s: %w", tenderId, err)
	}
	return reviews, nil
}
