package api

import (
	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/bll"
)

// Group ..
type Group struct {
	blls *bll.Blls
}

// CheckExists ..
func (g *Group) CheckExists(ctx *gear.Context) error {
	return nil
}

// BatchAdd ..
func (g *Group) BatchAdd(ctx *gear.Context) error {
	return nil
}

// Update ..
func (g *Group) Update(ctx *gear.Context) error {
	return nil
}

// Delete ..
func (g *Group) Delete(ctx *gear.Context) error {
	return nil
}

// BatchAddMembers ..
func (g *Group) BatchAddMembers(ctx *gear.Context) error {
	return nil
}

// RemoveMembers ..
func (g *Group) RemoveMembers(ctx *gear.Context) error {
	return nil
}

// GetLables ..
func (g *Group) GetLables(ctx *gear.Context) error {
	return nil
}

// GetSettings ..
func (g *Group) GetSettings(ctx *gear.Context) error {
	return nil
}

// RemoveLable ..
func (g *Group) RemoveLable(ctx *gear.Context) error {
	return nil
}

// UpdateSetting ..
func (g *Group) UpdateSetting(ctx *gear.Context) error {
	return nil
}

// RemoveSetting ..
func (g *Group) RemoveSetting(ctx *gear.Context) error {
	return nil
}
