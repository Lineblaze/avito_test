package postgresql

import (
	"fmt"
	openapi "github.com/Lineblaze/avito_gen"
	"github.com/google/uuid"
	"time"
	"zadanie-6105/backend/internal/domain"
	"zadanie-6105/backend/pkg/storage/postgres"
)

//go:generate ifacemaker -f postgres.go -o ../repository.go -i Repository -s PostgresRepository -p internal -y "Controller describes methods, implemented by the repository package."
type PostgresRepository struct {
	db postgres.Postgres
}

func NewPostgresRepository(db postgres.Postgres) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Employee

func (p *PostgresRepository) GetEmployeeByID(id int64) (*domain.Employee, error) {
	var employee domain.Employee

	err := p.db.QueryRow(`
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee
		WHERE id = $1`, id,
	).Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying employee: %w", err)
	}

	return &employee, nil
}

func (p *PostgresRepository) GetEmployeeByUsername(username string) (*domain.Employee, error) {
	var employee domain.Employee

	err := p.db.QueryRow(`
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee
		WHERE username = $1`, username,
	).Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying employee: %w", err)
	}

	return &employee, nil
}

func (p *PostgresRepository) CreateEmployee(employee *domain.Employee) (*domain.Employee, error) {
	var createdEmployee domain.Employee
	var createdAt time.Time

	err := p.db.QueryRow(`
        INSERT INTO employee(username, first_name, last_name, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, username, first_name, last_name, created_at`,
		employee.Username,
		employee.FirstName,
		employee.LastName,
	).Scan(
		&createdEmployee.ID,
		&createdEmployee.Username,
		&createdEmployee.FirstName,
		&createdEmployee.LastName,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	createdEmployee.CreatedAt = createdAt.Format(time.RFC3339)

	return &createdEmployee, nil
}

// Organization

func (p *PostgresRepository) GetOrganizationByID(id int64) (*domain.Organization, error) {
	var organization domain.Organization

	err := p.db.QueryRow(`
		SELECT id, name, description, type, created_at
		FROM organization
		WHERE id = $1`, id,
	).Scan(&organization.ID, &organization.Name, &organization.Description, &organization.Type, &organization.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("querying organization: %w", err)
	}

	return &organization, nil
}

func (p *PostgresRepository) CreateOrganization(organization *domain.Organization) (*domain.Organization, error) {
	var createdOrganization domain.Organization
	var createdAt time.Time

	err := p.db.QueryRow(`
        INSERT INTO organization(name, description, type, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, name, description, type, created_at`,
		organization.Name,
		organization.Description,
		organization.Type,
	).Scan(
		&createdOrganization.ID,
		&createdOrganization.Name,
		&createdOrganization.Description,
		&createdOrganization.Type,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	createdOrganization.CreatedAt = createdAt.Format(time.RFC3339)

	return &createdOrganization, nil
}

func (p *PostgresRepository) AssignEmployeeToOrganization(orgResp *domain.OrganizationResponsible) (*domain.OrganizationResponsible, error) {
	var exists bool
	err := p.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM organization_responsible WHERE user_id = $1
        )`, orgResp.UserID).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to check if user is already responsible for an organization: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user is already responsible for another organization")
	}

	var assign domain.OrganizationResponsible
	err = p.db.QueryRow(`
        INSERT INTO organization_responsible (organization_id, user_id)
        VALUES ($1, $2)
        RETURNING id, organization_id, user_id`,
		orgResp.OrganizationID,
		orgResp.UserID,
	).Scan(
		&assign.ID,
		&assign.OrganizationID,
		&assign.UserID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to assign employee to organization: %w", err)
	}

	return &assign, nil
}

func (p *PostgresRepository) IsUserResponsibleForOrganization(organizationId string) (bool, error) {
	var count int
	err := p.db.QueryRow(`
        SELECT COUNT(1) 
        FROM organization_responsible 
        WHERE organization_id = $1
    `, organizationId).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check user responsibility: %w", err)
	}

	return count > 0, nil
}

// Tenders

func (p *PostgresRepository) GetTenderByID(tenderId string) (*openapi.Tender, error) {
	var tender openapi.Tender
	var createdAt time.Time

	err := p.db.QueryRow(`
		SELECT id, name, description, service_type, status, organization_id, version, created_at
		FROM tender
		WHERE id = $1`, tenderId,
	).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying tender: %w", err)
	}

	tender.CreatedAt = createdAt.Format(time.RFC3339)

	return &tender, nil
}

func (p *PostgresRepository) GetTenders() ([]openapi.Tender, error) {
	rows, err := p.db.Query(`
        SELECT id, name, description, service_type, status, organization_id, version, created_at
        FROM tender
    `)
	if err != nil {
		return nil, fmt.Errorf("querying tenders: %w", err)
	}
	defer rows.Close()

	var tenders []openapi.Tender
	for rows.Next() {
		var tender openapi.Tender
		var createdAt time.Time
		if err = rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationId,
			&tender.Version,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("scanning tender: %w", err)
		}
		tender.CreatedAt = createdAt.Format(time.RFC3339)
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (p *PostgresRepository) GetUserTenders(userName string) ([]*openapi.Tender, error) {
	var tenders []*openapi.Tender
	var organizationId uuid.UUID

	err := p.db.QueryRow(`
        SELECT o.id 
        FROM employee e
        JOIN organization_responsible orp ON e.id = orp.user_id
        JOIN organization o ON orp.organization_id = o.id
        WHERE e.username = $1`, userName).Scan(&organizationId)

	if err != nil {
		return nil, fmt.Errorf("error fetching organization ID: %w", err)
	}

	rows, err := p.db.Query(`
        SELECT id, name, description, service_type, status, organization_id, version, created_at
        FROM tender
        WHERE organization_id = $1`, organizationId)

	if err != nil {
		return nil, fmt.Errorf("error fetching tenders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tender openapi.Tender
		var createdAt time.Time
		if err = rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationId,
			&tender.Version,
			&createdAt,
		); err != nil {
			return nil, err
		}
		tender.CreatedAt = createdAt.Format(time.RFC3339)
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (p *PostgresRepository) GetTenderByVersion(tenderID string, version string) (*openapi.Tender, error) {
	var tender openapi.Tender
	var createdAt time.Time

	err := p.db.QueryRow(`
	SELECT tender_id, name, description, service_type, status, organization_id, version, created_at
	FROM tender_version
	WHERE tender_id = $1 AND version = $2`,
		tenderID, version).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving tender by version: %w", err)
	}

	tender.CreatedAt = createdAt.Format(time.RFC3339)
	return &tender, nil
}

func (p *PostgresRepository) GetTenderStatus(tenderID string) (string, error) {
	var status string
	err := p.db.QueryRow(`
		SELECT status 
		FROM tender 
		WHERE id = $1`, tenderID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("failed to get tender status: %w", err)
	}
	return status, nil
}

func (p *PostgresRepository) CreateTender(tender *openapi.Tender) (*openapi.Tender, error) {
	var createdTender openapi.Tender
	var createdAt time.Time

	err := p.db.QueryRow(`
        INSERT INTO tender (name, description, service_type, status, organization_id, version, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        RETURNING id, name, description, service_type, status, organization_id, version, created_at`,
		tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.Version,
	).Scan(
		&createdTender.Id,
		&createdTender.Name,
		&createdTender.Description,
		&createdTender.ServiceType,
		&createdTender.Status,
		&createdTender.OrganizationId,
		&createdTender.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create tender: %w", err)
	}

	createdTender.CreatedAt = createdAt.Format(time.RFC3339)

	return &createdTender, nil
}

func (p *PostgresRepository) EditTender(tender *openapi.Tender) (*openapi.Tender, error) {
	_, err := p.db.Exec(`
		INSERT INTO tender_version (id, tender_id, name, description, service_type, status, organization_id, version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		uuid.New(), tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.Version, tender.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error saving tender version: %w", err)
	}

	var updatedTender openapi.Tender
	var createdAt time.Time

	err = p.db.QueryRow(`
		UPDATE tender
		SET name = $1, description = $2, service_type = $3, version = version + 1
		WHERE id = $4
		RETURNING id, name, description, service_type, status, organization_id, version, created_at`,
		tender.Name, tender.Description, tender.ServiceType, tender.Id).Scan(
		&updatedTender.Id,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.ServiceType,
		&updatedTender.Status,
		&updatedTender.OrganizationId,
		&updatedTender.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating tender: %w", err)
	}

	updatedTender.CreatedAt = createdAt.Format(time.RFC3339)

	return &updatedTender, nil
}

func (p *PostgresRepository) UpdateTenderStatus(tenderID string, status string) error {
	_, err := p.db.Exec(`
		UPDATE tender
		SET status = $1
		WHERE id = $2`, status, tenderID)
	if err != nil {
		return fmt.Errorf("failed to update tender status: %w", err)
	}
	return nil
}

// Bids

func (p *PostgresRepository) GetBidByID(bidId string) (*openapi.Bid, error) {
	var bid openapi.Bid
	var createdAt time.Time

	err := p.db.QueryRow(`
		SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
		FROM bid
		WHERE id = $1`, bidId,
	).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying bid: %w", err)
	}

	bid.CreatedAt = createdAt.Format(time.RFC3339)

	return &bid, nil
}

func (p *PostgresRepository) BidExistsByTenderID(tenderId string) (bool, error) {
	var count int

	err := p.db.QueryRow(`
        SELECT COUNT(*)
        FROM bid
        WHERE tender_id = $1`, tenderId,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check bid existence by tender ID: %w", err)
	}

	return count > 0, nil
}

func (p *PostgresRepository) GetBidsByTenderID(tenderId string) ([]*openapi.Bid, error) {
	var bids []*openapi.Bid

	rows, err := p.db.Query(`
        SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
        FROM bid
        WHERE tender_id = $1`, tenderId)

	if err != nil {
		return nil, fmt.Errorf("error fetching bids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bid openapi.Bid
		var createdAt time.Time
		if err = rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&createdAt,
		); err != nil {
			return nil, err
		}
		bid.CreatedAt = createdAt.Format(time.RFC3339)
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (p *PostgresRepository) GetUserBids(userName string) ([]*openapi.Bid, error) {
	var bids []*openapi.Bid
	var organizationId uuid.UUID

	err := p.db.QueryRow(`
        SELECT o.id 
        FROM employee e
        JOIN organization_responsible orp ON e.id = orp.user_id
        JOIN organization o ON orp.organization_id = o.id
        WHERE e.username = $1`, userName).Scan(&organizationId)

	if err != nil {
		return nil, fmt.Errorf("error fetching organization ID: %w", err)
	}

	rows, err := p.db.Query(`
        SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
        FROM bid
        WHERE author_id = $1`, organizationId)

	if err != nil {
		return nil, fmt.Errorf("error fetching tenders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bid openapi.Bid
		var createdAt time.Time
		if err = rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorId,
			&bid.AuthorType,
			&bid.Version,
			&createdAt,
		); err != nil {
			return nil, err
		}
		bid.CreatedAt = createdAt.Format(time.RFC3339)
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (p *PostgresRepository) GetBidByVersion(bidID string, version string) (*openapi.Bid, error) {
	var bid openapi.Bid
	var createdAt time.Time

	err := p.db.QueryRow(`
	SELECT bid_id, name, description, status, tender_id, author_type, author_id, version, created_at
	FROM bid_version
	WHERE bid_id = $1 AND version = $2`,
		bidID, version).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving bid by version: %w", err)
	}

	bid.CreatedAt = createdAt.Format(time.RFC3339)
	return &bid, nil
}

func (p *PostgresRepository) GetBidStatus(bidID string) (string, error) {
	var status string
	err := p.db.QueryRow(`
		SELECT status 
		FROM bid
		WHERE id = $1`, bidID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("failed to get bid status: %w", err)
	}
	return status, nil
}

func (p *PostgresRepository) CreateBid(bid *openapi.Bid) (*openapi.Bid, error) {
	var createdBid openapi.Bid
	var createdAt time.Time

	err := p.db.QueryRow(`
        INSERT INTO bid (name, description, status, tender_id, author_id, author_type, version, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING id, name, description, status, tender_id, author_id, author_type, version, created_at`,
		bid.Name, bid.Description, bid.Status, bid.TenderId, bid.AuthorId, bid.AuthorType, bid.Version,
	).Scan(
		&createdBid.Id,
		&createdBid.Name,
		&createdBid.Description,
		&createdBid.Status,
		&createdBid.TenderId,
		&createdBid.AuthorId,
		&createdBid.AuthorType,
		&createdBid.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	createdBid.CreatedAt = createdAt.Format(time.RFC3339)

	return &createdBid, nil
}

func (p *PostgresRepository) EditBid(bid *openapi.Bid) (*openapi.Bid, error) {
	_, err := p.db.Exec(`
		INSERT INTO bid_version (id, bid_id, name, description, status, tender_id, author_type, author_id, version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		uuid.New(), bid.Id, bid.Name, bid.Description, bid.Status, bid.TenderId, bid.AuthorType, bid.AuthorId, bid.Version, bid.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error saving bid version: %w", err)
	}

	var updatedBid openapi.Bid
	var createdAt time.Time

	err = p.db.QueryRow(`
		UPDATE bid
		SET name = $1, description = $2, status = $3, version = version + 1
		WHERE id = $4
		RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at`,
		bid.Name, bid.Description, bid.Status, bid.Id).Scan(
		&updatedBid.Id,
		&updatedBid.Name,
		&updatedBid.Description,
		&updatedBid.Status,
		&updatedBid.TenderId,
		&updatedBid.AuthorType,
		&updatedBid.AuthorId,
		&updatedBid.Version,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating bid: %w", err)
	}

	updatedBid.CreatedAt = createdAt.Format(time.RFC3339)

	return &updatedBid, nil
}

func (p *PostgresRepository) UpdateBidStatus(bidID string, status string) error {
	_, err := p.db.Exec(`
		UPDATE bid
		SET status = $1
		WHERE id = $2`, status, bidID)
	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}
	return nil
}
