package http

import (
	openapi "github.com/Lineblaze/avito_gen"
	"github.com/gofiber/fiber/v3"
	useCase "zadanie-6105/backend/internal"
	"zadanie-6105/backend/internal/domain"
	"zadanie-6105/backend/pkg/logger"
)

//go:generate ifacemaker -f handler.go -o ../../handler.go -i Handler -s Handler -p internal -y "Controller describes methods, implemented by the http package."
type Handler struct {
	useCase useCase.UseCase
	logger  *logger.ApiLogger
}

func NewHandler(useCase useCase.UseCase, logger *logger.ApiLogger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Ping

func (h Handler) Ping() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	}
}

// Employees

func (h Handler) CreateEmployee() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.CreateEmployeeRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		employee, err := h.useCase.CreateEmployee(&req)
		if err != nil {
			h.logger.Errorf("Failed to create employee", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(employee)
	}
}

// Organizations

func (h Handler) CreateOrganization() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.CreateOrganizationRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		organization, err := h.useCase.CreateOrganization(&req)
		if err != nil {
			h.logger.Errorf("Failed to create flat", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(organization)
	}
}

func (h Handler) AssignEmployeeToOrganization() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req domain.AssignEmployeeToOrganizationRequest

		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		orgResp, err := h.useCase.AssignEmployeeToOrganization(&req)
		if err != nil {
			h.logger.Errorf("Failed to assign employee to organization: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(orgResp)
	}
}

// Tenders

func (h Handler) GetTenders() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenders, err := h.useCase.GetTenders()
		if err != nil {
			h.logger.Errorf("Failed to get tenders", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(tenders)
	}
}

func (h *Handler) GetUserTenders() fiber.Handler {
	return func(c fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username is required"})
		}
		tenders, err := h.useCase.GetUserTenders(username)
		if err != nil {
			h.logger.Errorf("Failed to fetch tenders", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(tenders)
	}
}

func (h *Handler) GetTenderStatus() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderID := c.Params("tenderId")

		status, err := h.useCase.GetTenderStatus(tenderID)
		if err != nil {
			h.logger.Errorf("Failed to get tender status: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": status})
	}
}

func (h Handler) CreateTender() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req openapi.CreateTenderRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		isResponsible, err := h.useCase.CheckUserOrganizationResponsibility(req.OrganizationId)
		if err != nil {
			h.logger.Errorf("Failed to check user responsibility", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		if !isResponsible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not responsible for this organization"})
		}

		tender, err := h.useCase.CreateTender(&req)
		if err != nil {
			h.logger.Errorf("Failed to create tender", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(tender)
	}
}

func (h *Handler) EditTender() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderID := c.Params("tenderId")
		var req openapi.EditTenderRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		updatedTender, err := h.useCase.EditTender(tenderID, &req)
		if err != nil {
			h.logger.Errorf("Failed to update tender", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(updatedTender)
	}
}

func (h *Handler) UpdateTenderStatus() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderID := c.Params("tenderId")
		status := c.Params("status")

		err := h.useCase.UpdateTenderStatus(tenderID, status)
		if err != nil {
			h.logger.Errorf("Failed to update tender status", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *Handler) RollbackTender() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderID := c.Params("tenderId")
		version := c.Params("version")

		rolledBackTender, err := h.useCase.RollbackTender(tenderID, version)
		if err != nil {
			h.logger.Errorf("Failed to rollback tender", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(rolledBackTender)
	}
}

// Bids

func (h *Handler) GetUserBids() fiber.Handler {
	return func(c fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username is required"})
		}
		bids, err := h.useCase.GetUserBids(username)
		if err != nil {
			h.logger.Errorf("Failed to fetch tenders", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(bids)
	}
}

func (h *Handler) GetBidsByTenderID() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderID := c.Params("tenderId")
		if tenderID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tenderId is required"})
		}

		bids, err := h.useCase.GetBidsByTenderID(tenderID)
		if err != nil {
			h.logger.Errorf("Failed to get bids by tender ID: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(bids)
	}
}

func (h *Handler) GetBidStatus() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidID := c.Params("bidId")
		if bidID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bidId is required"})
		}

		status, err := h.useCase.GetBidStatus(bidID)
		if err != nil {
			h.logger.Errorf("Failed to get bid status: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": status})
	}
}

func (h Handler) CreateBid() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req openapi.CreateBidRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		isResponsible, err := h.useCase.CheckUserOrganizationResponsibility(req.OrganizationId)
		if err != nil {
			h.logger.Errorf("Failed to check user responsibility", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		if !isResponsible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not responsible for this organization"})
		}

		bid, err := h.useCase.CreateBid(&req)
		if err != nil {
			h.logger.Errorf("Failed to create tender", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(bid)
	}
}

func (h *Handler) EditBid() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidID := c.Params("bidId")
		if bidID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bidId is required"})
		}

		var req openapi.EditBidRequest
		if err := c.Bind().Body(&req); err != nil {
			h.logger.Errorf("Failed to parse request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		updatedBid, err := h.useCase.EditBid(bidID, &req)
		if err != nil {
			h.logger.Errorf("Failed to edit bid: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(updatedBid)
	}
}

func (h *Handler) UpdateBidStatus() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidID := c.Params("bidId")
		status := c.Params("status")
		if bidID == "" || status == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bidId and status are required"})
		}

		err := h.useCase.UpdateBidStatus(bidID, status)
		if err != nil {
			h.logger.Errorf("Failed to update bid status: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *Handler) RollbackBid() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidID := c.Params("bidId")
		version := c.Params("version")

		rolledBackBid, err := h.useCase.RollbackBid(bidID, version)
		if err != nil {
			h.logger.Errorf("Failed to rollback bid", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(rolledBackBid)
	}
}

func (h Handler) SubmitBidDecision() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidId := c.Params("bidId")
		decision := c.Params("decision")
		username := c.Params("username")

		if decision != "Approved" && decision != "Rejected" {
			h.logger.Errorf("Invalid decision")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid decision"})
		}

		isResponsible, err := h.useCase.CheckUserOrganizationResponsibilityByUsername(username)
		if err != nil {
			h.logger.Errorf("Failed to check user responsibility", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		if !isResponsible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not responsible for this organization"})
		}

		err = h.useCase.SubmitBidDecision(bidId, decision, username)
		if err != nil {
			h.logger.Errorf("Failed to submit bid decision", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Bid decision submitted successfully"})
	}
}

func (h Handler) SubmitBidFeedback() fiber.Handler {
	return func(c fiber.Ctx) error {
		bidId := c.Params("bidId")
		feedback := c.Params("feedback")
		username := c.Params("username")

		if len(feedback) > 1000 {
			h.logger.Errorf("Feedback is too long")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Feedback exceeds character limit"})
		}

		isResponsible, err := h.useCase.CheckUserOrganizationResponsibilityByUsername(username)
		if err != nil {
			h.logger.Errorf("Failed to check user responsibility", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		if !isResponsible {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not responsible for this organization"})
		}

		err = h.useCase.SubmitBidFeedback(bidId, feedback, username)
		if err != nil {
			h.logger.Errorf("Failed to submit bid feedback", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Feedback submitted successfully"})
	}
}

func (h Handler) GetBidReviews() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenderId := c.Params("tenderId")
		username := c.Params("username")

		isResponsible, err := h.useCase.CheckUserOrganizationResponsibilityByUsername(username)
		if err != nil {
			h.logger.Errorf("Failed to check user responsibility: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		if !isResponsible {
			h.logger.Warnf("User %s is not responsible for tender %s", username, tenderId)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User not responsible for this organization"})
		}

		reviews, err := h.useCase.GetBidReviewsByTenderId(tenderId)
		if err != nil {
			h.logger.Errorf("Failed to get bid reviews: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "InternalServerError"})
		}

		return c.Status(fiber.StatusOK).JSON(reviews)
	}
}
