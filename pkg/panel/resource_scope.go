package panel

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

var defaultInternalResourceSlugs = map[string]struct{}{
	"users":         {},
	"accounts":      {},
	"sessions":      {},
	"verifications": {},
}

type registrySnapshot struct {
	resources             map[string]resource.Resource
	publicResources       map[string]resource.Resource
	internalResourceSlugs map[string]struct{}
	pages                 map[string]page.Page
}

func (p *Panel) registerResourceWithScope(slug string, res resource.Resource, internal bool) {
	if p == nil || res == nil || slug == "" {
		return
	}

	p.registryMu.Lock()
	defer p.registryMu.Unlock()

	if p.resources == nil {
		p.resources = make(map[string]resource.Resource)
	}
	if p.publicResources == nil {
		p.publicResources = make(map[string]resource.Resource)
	}
	if p.internalResourceSlugs == nil {
		p.internalResourceSlugs = make(map[string]struct{})
	}

	if res.GetDialogType() == "" {
		res.SetDialogType(resource.DialogTypeSheet)
	}
	p.resources[slug] = res

	if internal {
		p.internalResourceSlugs[slug] = struct{}{}
		delete(p.publicResources, slug)
		p.publishRegistrySnapshotLocked()
		return
	}

	if _, ok := defaultInternalResourceSlugs[slug]; ok {
		p.publishRegistrySnapshotLocked()
		return
	}
	if _, ok := p.internalResourceSlugs[slug]; ok {
		p.publishRegistrySnapshotLocked()
		return
	}

	p.publicResources[slug] = res
	p.publishRegistrySnapshotLocked()
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

	snapshot := p.loadRegistrySnapshot()
	if snapshot == nil {
		return false
	}

	_, ok := snapshot.internalResourceSlugs[slug]
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

	snapshot := p.loadRegistrySnapshot()
	if snapshot == nil {
		return nil, false
	}

	res, ok := snapshot.resources[slug]
	if !ok {
		return nil, false
	}

	if !p.isResourceAccessibleForRequest(c, slug) {
		return nil, false
	}

	return res, true
}

func (p *Panel) loadRegistrySnapshot() *registrySnapshot {
	if p == nil {
		return nil
	}

	if raw := p.registrySnapshot.Load(); raw != nil {
		if snapshot, ok := raw.(*registrySnapshot); ok {
			return snapshot
		}
	}

	p.registryMu.RLock()
	defer p.registryMu.RUnlock()

	return p.buildRegistrySnapshotLocked()
}

func (p *Panel) buildRegistrySnapshotLocked() *registrySnapshot {
	if p == nil {
		return nil
	}

	return &registrySnapshot{
		resources:             cloneResourceMap(p.resources),
		publicResources:       cloneResourceMap(p.publicResources),
		internalResourceSlugs: cloneInternalSlugMap(p.internalResourceSlugs),
		pages:                 clonePageMap(p.pages),
	}
}

func (p *Panel) publishRegistrySnapshotLocked() {
	if p == nil {
		return
	}
	p.registrySnapshot.Store(p.buildRegistrySnapshotLocked())
}

func (p *Panel) pagesSnapshot() map[string]page.Page {
	snapshot := p.loadRegistrySnapshot()
	if snapshot == nil {
		return map[string]page.Page{}
	}
	return snapshot.pages
}

func cloneResourceMap(src map[string]resource.Resource) map[string]resource.Resource {
	if len(src) == 0 {
		return map[string]resource.Resource{}
	}
	dst := make(map[string]resource.Resource, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func clonePageMap(src map[string]page.Page) map[string]page.Page {
	if len(src) == 0 {
		return map[string]page.Page{}
	}
	dst := make(map[string]page.Page, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func cloneInternalSlugMap(src map[string]struct{}) map[string]struct{} {
	if len(src) == 0 {
		return map[string]struct{}{}
	}
	dst := make(map[string]struct{}, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
