package calculator

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"CalculatorWebService/calculator/storage"
)

type Request struct {
	Operand1 float64 `form:"operand1" json:"operand1"`
	Operand2 float64 `form:"operand2" json:"operand2"`
}

type Response struct {
	Result     float64 `json:"result"`
	Operation  string  `json:"operation"`
	Expression string  `json:"expression"`
}

type RecentResponse struct {
	Calculations []string `json:"calculations"`
}

type Handler struct {
	Storage storage.Storage
}

// Considering the scope of the service it's okay to use a single instance of Handler and perform the logic inside methods.
// In real world scenarios, it's better to separate business logic from handlers.
// handlers should ideally just handle HTTP specifics (parsing requests, validating requests, forming responses) and delegate business logic to separate services.
// For example, we could have a CalculatorService struct that would handle the operations and storage interactions.
// Handlers would then call methods on that service.

func NewCalculationHandler(storage storage.Storage) *Handler {
	return &Handler{
		Storage: storage,
	}
}

// Addition handler
func (h *Handler) Addition(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := req.Operand1 + req.Operand2
	expression := formatExpression(req.Operand1, req.Operand2, "+", result)

	h.Storage.Store(expression)

	response := Response{
		Result:     result,
		Operation:  "addition",
		Expression: expression,
	}

	c.JSON(http.StatusOK, response)
}

// Subtraction handler
func (h *Handler) Subtraction(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := req.Operand1 - req.Operand2
	expression := formatExpression(req.Operand1, req.Operand2, "-", result)

	h.Storage.Store(expression)

	response := Response{
		Result:     result,
		Operation:  "subtraction",
		Expression: expression,
	}

	c.JSON(http.StatusOK, response)
}

// Multiplication handler with error handling
func (h *Handler) Multiplication(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := req.Operand1 * req.Operand2
	expression := formatExpression(req.Operand1, req.Operand2, "*", result)

	h.Storage.Store(expression)

	response := Response{
		Result:     result,
		Operation:  "multiplication",
		Expression: expression,
	}

	c.JSON(http.StatusOK, response)
}

// Division handler with error handling
func (h *Handler) Division(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Operand2 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Division by zero is not allowed"})
		return
	}

	result := req.Operand1 / req.Operand2
	expression := formatExpression(req.Operand1, req.Operand2, "/", result)

	h.Storage.Store(expression)

	response := Response{
		Result:     result,
		Operation:  "division",
		Expression: expression,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetRecentCalculations(c *gin.Context) {
	n := 5 // default
	if nStr := c.Param("n"); nStr != "" {
		if parsed, err := strconv.Atoi(nStr); err == nil && parsed > 0 && parsed <= 20 {
			n = parsed
		}
	}

	calculations := h.Storage.GetRecent(n)

	response := RecentResponse{
		Calculations: calculations,
	}

	c.JSON(http.StatusOK, response)
}

func formatExpression(a, b float64, operator string, result float64) string {
	return strconv.FormatFloat(a, 'f', -1, 64) + " " + operator + " " +
		strconv.FormatFloat(b, 'f', -1, 64) + " = " +
		strconv.FormatFloat(result, 'f', -1, 64)
}
