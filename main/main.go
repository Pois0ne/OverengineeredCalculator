package main

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Expression represents an arithmetic expression
type Expression struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Result    float64   `json:"result"`
	Processed bool      `json:"processed"`
}

// Agent represents Computational Agent routine and status
type Agent struct {
	Status string `json:"status"`
	Num    int    `json:"num"`
}

var (
	expressions     = make(map[string]*Expression)
	agents          = []*Agent{}
	timer           time.Duration
	expressionMutex sync.Mutex
)

func main() {

	timer = 10 * time.Second //Fake process time. Default : 10 sec

	http.HandleFunc("/expression", handleExpression)                 //Add expression to compute
	http.HandleFunc("/expressions", handleExpressions)               //Show list of expressions
	http.HandleFunc("/timer", handleTimer)                           //Change process fake timer
	http.HandleFunc("/computationalAgent", handleComputationalAgent) //add agents
	http.HandleFunc("/agentsList", handleAgentsList)                 //shows agents and their statuses

	go startComputationalAgents(5)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAgentsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(agents)
}

func handleComputationalAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	n, err := strconv.Atoi(r.FormValue("add"))

	if err != nil {
		http.Error(w, "Value must be int", http.StatusBadRequest)
		return
	} else if n <= 0 {
		http.Error(w, "Value must be grater than 0", http.StatusBadRequest)
		return
	}

	log.Println("staring", n, "new Computational Agents ")

	startComputationalAgents(n)

}

func handleTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	duration, err := strconv.Atoi(r.FormValue("TimeToProceed"))

	if err != nil {
		http.Error(w, "Wrong time format", http.StatusBadRequest)
		return
	} else if duration <= 0 {
		http.Error(w, "Time must be grater than 0", http.StatusBadRequest)
		return
	}

	log.Println("Task process time has been changed from", timer, "to:", duration, "seconds")
	timer = time.Second * time.Duration(duration)
	log.Println("Task process time now =", timer)

}

func handleExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	content := r.FormValue("expression")
	content = strings.ReplaceAll(content, "p", "+") //Костыль с "+". Не смог понять как это исправить адекватным методом
	log.Println(content)

	if content == "" {
		http.Error(w, "Expression is required", http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		id = generateID()
	}
	expression := &Expression{
		ID:        id,
		Content:   content,
		Status:    "processing",
		Created:   time.Now(),
		Updated:   time.Now(),
		Processed: false,
	}
	expressionMutex.Lock()
	expressions[id] = expression
	expressionMutex.Unlock()
	w.WriteHeader(http.StatusOK)
}
func handleExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	expressionMutex.Lock()
	defer expressionMutex.Unlock()
	expressionsList := make([]*Expression, 0, len(expressions))
	for _, expression := range expressions {
		expressionsList = append(expressionsList, expression)
	}
	json.NewEncoder(w).Encode(expressionsList)
}

func startComputationalAgents(numAgents int) {

	log.Println("Starting", numAgents, "Computational Agents...")

	for i := 0; i < numAgents; i++ {
		agent := &Agent{
			Status: "idle",
			Num:    len(agents) + 1,
		}

		agents = append(agents, agent)
		time.Sleep(1 * time.Millisecond)
		go agent.computationalAgent()
	}

}
func (a *Agent) computationalAgent() {

	log.Println("Successfully started Computational agent #", a.Num)

	for {
		time.Sleep(1 * time.Second) //Agent expression check delay
		expression := getNextExpression()
		if expression == nil {
			a.Status = "idle"
			continue
		}
		if !expression.Processed {
			a.Status = "working"
			result, err := calculateExpression(expression.Content)
			updateExpressionResult(expression, result, err)
		}
	}
}
func getNextExpression() *Expression {
	expressionMutex.Lock()
	defer expressionMutex.Unlock()
	for _, expression := range expressions {
		if expression.Status == "processing" {
			expression.Status = "in progress"
			expression.Updated = time.Now()
			return expression
		}
	}
	return nil
}
func calculateExpression(expression string) (float64, error) {
	time.Sleep(timer) //Fake delay

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}
func updateExpressionResult(expression *Expression, result float64, err error) {
	expressionMutex.Lock()
	defer expressionMutex.Unlock()

	switch {
	case err == nil:
		expression.Status = "completed"
		expression.Updated = time.Now()
		expression.Result = result
		expression.Processed = true
		log.Println("\n\n##################\n#Successfully completed expression \n#ID:", expression.ID, "\n#Content:", expression.Content, "\n#Result:", expression.Result, "\n##################\n\n")

	case err != nil:
		log.Println("Fail while parsing expression:", err)

		expression.Status = "invalid expression"
		expression.Updated = time.Now()
		expression.Result = result
		expression.Processed = true
	}
}
func generateID() string {
	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprintf("%d", rand.Intn(10000))
	return id
}
