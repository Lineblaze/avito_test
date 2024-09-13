https://github.com/Lineblaze/avito_gen - Пакет c OpenAPI генерацией

Используются следующие роуты :

	r.Get(`/ping`, h.Ping())

	r.Post(`/employee/new`, h.CreateEmployee())
	r.Post(`/organization/new`, h.CreateOrganization())
	r.Post(`/assign`, h.AssignEmployeeToOrganization())

	r.Get(`/tenders`, h.GetTenders())
	r.Get(`/tenders/my`, h.GetUserTenders())
	r.Get("/tenders/:tenderId/status", h.GetTenderStatus())
	r.Post(`/tenders/new`, h.CreateTender())
	r.Patch(`/tenders/:tenderId/edit`, h.EditTender())
	r.Put(`/tenders/:tenderId/rollback/:version`, h.RollbackTender())
	r.Put("/tenders/:tenderId/status/:status", h.UpdateTenderStatus())

	r.Get(`/bids/my`, h.GetUserBids())
	r.Get(`/bids/:tenderId/list`, h.GetBidsByTenderID())
	r.Get("/bids/:bidId/status", h.GetBidStatus())
	r.Post(`/bids/new`, h.CreateBid())
	r.Patch(`/bids/:bidId/edit`, h.EditBid())
	r.Put(`/bids/:bidId/rollback/:version`, h.RollbackBid())
	r.Put("/bids/:bidId/status/:status", h.UpdateBidStatus())

	r.Get("/bids/:tenderId/reviews/:username", h.GetBidReviews())
	r.Put("/bids/:bidId/submit_decision/:decision/:username", h.SubmitBidDecision())
	r.Put("/bids/:bidId/feedback/:feedback/:username", h.SubmitBidFeedback())

Есть Avito.postman_collection в файлах проекта