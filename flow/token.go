package flow

import (
	"encoding/json"

	"github.com/petri-nets/workflow/wfmod"
)

// Token token
type Token struct {
	WfToken   wfmod.WfToken
	WfContext wfmod.WfContextType
}

// NewToken new token
func NewToken(cs *Case, context string, placeID int) Token {
	wfToken := wfmod.WfToken{
		AppID:      cs.WfCase.AppID,
		WorkflowID: cs.WfCase.WorkflowID,
		CaseID:     cs.WfCase.ID,
		PlaceID:    placeID,
		Context:    context,
		Status:     wfmod.TokenStatusFree,
	}

	flowDao.SaveToken(&wfToken)
	return buildToken(wfToken)
}

// GetToken get exist token
func GetToken(caseID int) Token {
	// 通过caseID查找token
	return Token{
		WfToken: wfmod.WfToken{},
	}
}

func buildToken(wfToken wfmod.WfToken) Token {
	token := Token{
		WfToken:   wfToken,
		WfContext: wfmod.WfContextType{},
	}

	json.Unmarshal([]byte(wfToken.Context), &token.WfContext)

	return token
}

// GetPlacesTokens get places tokens
func GetPlacesTokens(cs *Case, placeIDList []int, status wfmod.TokenStatusType) []Token {
	wfTokens := flowDao.GetTokensByPlaces(&cs.WfCase, placeIDList, status)

	var tokens []Token
	for _, wfToken := range wfTokens {
		tokens = append(tokens, buildToken(wfToken))
	}
	return tokens
}

// MergeTokens merge tokens
func MergeTokens(tokens []Token) Token {
	var tk Token
	if len(tokens) > 0 {
		tk = tokens[0]
	}

	for _, item := range tokens {
		if tk.WfContext == nil {
			tk.WfContext = wfmod.WfContextType{}
		}

		if item.WfContext != nil {
			for k, v := range item.WfContext {
				tk.WfContext[k] = v
			}
		}
	}
	return tk
}

// IsFree check token has consumed
func (t *Token) IsFree() bool {
	// move token
	wfTokens := flowDao.GetTokenByIDList(t.WfToken.AppID, []int{t.WfToken.ID})
	if len(wfTokens) > 0 {
		t.WfToken = wfTokens[0]
	}
	return t.WfToken.Status != wfmod.TokenStatusFree
}

// Lock lock the token
func (t *Token) Lock() error {
	t.WfToken.Status = wfmod.TokenStatusLock
	flowDao.SaveToken(&t.WfToken)
	return nil
}

// Consume consume the token
func (t *Token) Consume() error {
	t.WfToken.Status = wfmod.TokenStatusConsume
	flowDao.SaveToken(&t.WfToken)
	return nil
}
