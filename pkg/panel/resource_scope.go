package panel

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

var defaultInternalResourceSlugs = map[string]struct{}{
	"users":         {},
	"accounts":      {},
	"sessions":      {},
	"verifications": {},
}

func (p *Panel) registerResourceWithScope(slug string, res resource.Resource, internal bool) {
	if p == nil || res == nil || slug == "" {
		return
	}

	if p.resources == nil {
		p.resources = make(map[string]resource.Resource)
	}
	if p.publicResources == nil {
		p.publicResources = make(map[string]resource.Resource)
	}
	if p.internalResourceSlugs == nil {
		p.internalResourceSlugs = make(map[string]struct{})
	}

	res.SetDialogType(resource.DialogTypeSheet)
	p.resources[slug] = res

	if internal {
		p.internalResourceSlugs[slug] = struct{}{}
		delete(p.publicResources, slug)
		return
	}

	if p.isInternalResourceSlug(slug) {
		return
	}

	p.publicResources[slug] = res
}

func (p *Panel) registerSystemResource(res resource.Resource) {
	if res == nil {
		return
	}
	p.registerResourceWithScope(res.Slug(), res, true)
}

func (p *Panel) isInternalResourceSlug(slug string) bool {
	if slug == "" {
		return false
	}

	if _, ok := defaultInternalResourceSlugs[slug]; ok {
		return true
	}

	if p == nil {
		return false
	}
	_, ok := p.internalResourceSlugs[slug]
	return ok
}

func (p *Panel) isResourceAccessibleForRequest(c *context.Context, slug string) bool {
	if p == nil {
		return false
	}

	if c == nil {
		return true
	}

	if apiKeyAuth, ok := c.Locals(middleware.APIKeyAuthenticatedLocalKey).(bool); ok && apiKeyAuth {
		return !p.isInternalResourceSlug(slug)
	}

	return true
}

func (p *Panel) resolveResourceForRequest(c *context.Context, slug string) (resource.Resource, bool) {
	if p == nil || slug == "" {
		return nil, false
	}

	res, ok := p.resources[slug]
	if !ok {
		return nil, false
	}

	if !p.isResourceAccessibleForRequest(c, slug) {
		return nil, false
	}

	return res, true
}
