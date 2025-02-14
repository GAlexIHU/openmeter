package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/openmeterio/openmeter/api"
	"github.com/openmeterio/openmeter/internal/server/authenticator"
	"github.com/openmeterio/openmeter/pkg/contextx"
	"github.com/openmeterio/openmeter/pkg/models"
)

// CreatePortalToken creates a new portal token.
func (a *Router) CreatePortalToken(w http.ResponseWriter, r *http.Request) {
	ctx := contextx.WithAttr(r.Context(), "operation", "createPortalToken")

	if a.config.PortalTokenStrategy == nil {
		err := fmt.Errorf("not implemented: portal is not enabled")

		models.NewStatusProblem(ctx, err, http.StatusNotImplemented).Respond(w, r)

		return
	}

	// Parse request body
	body := &api.CreatePortalTokenJSONRequestBody{}
	if err := render.DecodeJSON(r.Body, body); err != nil {
		err := fmt.Errorf("decode json: %w", err)

		models.NewStatusProblem(ctx, err, http.StatusBadRequest).Respond(w, r)

		return
	}

	t, err := a.config.PortalTokenStrategy.Generate(body.Subject, body.AllowedMeterSlugs, body.ExpiresAt)
	if err != nil {
		err := fmt.Errorf("generate portal token: %w", err)

		a.config.ErrorHandler.HandleContext(ctx, err)
		models.NewStatusProblem(ctx, err, http.StatusInternalServerError).Respond(w, r)

		return
	}

	render.JSON(w, r, api.PortalToken{
		Id:                t.Id,
		Token:             t.Token,
		ExpiresAt:         t.ExpiresAt,
		Subject:           t.Subject,
		AllowedMeterSlugs: t.AllowedMeterSlugs,
	})
}

func (a *Router) ListPortalTokens(w http.ResponseWriter, r *http.Request, params api.ListPortalTokensParams) {
	ctx := contextx.WithAttr(r.Context(), "operation", "listPortalTokens")

	err := fmt.Errorf("not implemented: portal token listing is an OpenMeter Cloud only feature")

	models.NewStatusProblem(ctx, err, http.StatusNotImplemented).Respond(w, r)
}

func (a *Router) InvalidatePortalTokens(w http.ResponseWriter, r *http.Request) {
	ctx := contextx.WithAttr(r.Context(), "operation", "invalidatePortalTokens")

	err := fmt.Errorf("not implemented: portal token invalidation is an OpenMeter Cloud only feature")

	models.NewStatusProblem(ctx, err, http.StatusNotImplemented).Respond(w, r)
}

func (a *Router) QueryPortalMeter(w http.ResponseWriter, r *http.Request, meterSlug string, params api.QueryPortalMeterParams) {
	ctx := contextx.WithAttr(r.Context(), "operation", "queryPortalMeter")
	ctx = contextx.WithAttr(ctx, "meterSlug", meterSlug)
	ctx = contextx.WithAttr(ctx, "params", params) // TODO: we should probable NOT add this to the context

	subject := authenticator.GetAuthenticatedSubject(ctx)
	if subject == "" {
		err := fmt.Errorf("not authenticated")
		models.NewStatusProblem(ctx, err, http.StatusUnauthorized).Respond(w, r)
		return
	}

	a.QueryMeter(w, r, meterSlug, api.QueryMeterParams{
		From:           params.From,
		To:             params.To,
		WindowSize:     params.WindowSize,
		WindowTimeZone: params.WindowTimeZone,
		Subject:        &[]string{subject},
		GroupBy:        params.GroupBy,
	})
}
