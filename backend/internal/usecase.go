// Code generated by ifacemaker; DO NOT EDIT.

package internal

import (
	openapi "github.com/Lineblaze/avito_gen"
)

// Controller describes methods, implemented by the usecase package.
type UseCase interface {
	GetEmployeeByID(employeeID int64) (*openapi.Employee, error)
	GetEmployeeByUserName(username string) (*openapi.Employee, error)
	CreateEmployee(req *openapi.CreateEmployeeRequest) (*openapi.Employee, error)
	GetOrganizationByID(organizationID int64) (*openapi.Organization, error)
	CreateOrganization(req *openapi.CreateOrganizationRequest) (*openapi.Organization, error)
	AssignEmployeeToOrganization(req *openapi.AssignEmployeeToOrganizationRequest) (*openapi.OrganizationResponsible, error)
	CheckUserOrganizationResponsibility(organizationId string) (bool, error)
	CheckUserOrganizationResponsibilityByUsername(username string) (bool, error)
	GetTenders() ([]openapi.Tender, error)
	GetTenderByID(tenderID string) (*openapi.Tender, error)
	GetUserTenders(userName string) ([]*openapi.Tender, error)
	GetTenderStatus(tenderID string) (string, error)
	CreateTender(req *openapi.CreateTenderRequest) (*openapi.Tender, error)
	EditTender(tenderID string, req *openapi.EditTenderRequest) (*openapi.Tender, error)
	UpdateTenderStatus(tenderID string, status string) error
	RollbackTender(tenderID string, version string) (*openapi.Tender, error)
	CanUserAccessTender(tenderID string) (bool, error)
	GetBidByID(bidID string) (*openapi.Bid, error)
	GetUserBids(userName string) ([]*openapi.Bid, error)
	GetBidsByTenderID(tenderID string) ([]*openapi.Bid, error)
	GetBidStatus(bidID string) (string, error)
	CreateBid(req *openapi.CreateBidRequest) (*openapi.Bid, error)
	EditBid(bidID string, req *openapi.EditBidRequest) (*openapi.Bid, error)
	UpdateBidStatus(bidID string, status string) error
	RollbackBid(bidID string, version string) (*openapi.Bid, error)
	SubmitBidDecision(bidId string, decision string, username string) error
	IsTenderClosed(bidId string) (bool, error)
	SubmitBidFeedback(bidId string, feedback string, username string) error
	GetBidReviewsByTenderId(tenderId string) ([]openapi.BidReview, error)
}
