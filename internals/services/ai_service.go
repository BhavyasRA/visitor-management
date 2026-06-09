package services

type AIDocumentResponse struct {
	DocumentType   string
	DocumentNumber string
}

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

func (s *AIService) ExtractDocumentDetails(
	documentID uint,
	documentURL string,
) (*AIDocumentResponse, error) {

	// Later call AI team API here.
	// Send:
	// document_id = documentID
	// document_url = documentURL

	return &AIDocumentResponse{
		DocumentType:   "document_type",
		DocumentNumber: "document_number",
	}, nil
}
