// Code generated from the VEED OpenAPI spec. DO NOT EDIT.

package veed

type services struct {
	// Fabric accesses the fabric-1.0 model.
	Fabric *FabricService
	// Lipsync20 accesses the lipsync-2.0 model.
	Lipsync20 *Lipsync20Service
	// VideoBackgroundRemoval accesses the video-background-removal model.
	VideoBackgroundRemoval *VideoBackgroundRemovalService
	// VideoBackgroundRemovalFast accesses the video-background-removal-fast model.
	VideoBackgroundRemovalFast *VideoBackgroundRemovalFastService
	// VideoBackgroundRemovalGreenScreen accesses the video-background-removal-green-screen model.
	VideoBackgroundRemovalGreenScreen *VideoBackgroundRemovalGreenScreenService
}

func (c *Client) initServices() {
	c.Fabric = &FabricService{client: c}
	c.Lipsync20 = &Lipsync20Service{client: c}
	c.VideoBackgroundRemoval = &VideoBackgroundRemovalService{client: c}
	c.VideoBackgroundRemovalFast = &VideoBackgroundRemovalFastService{client: c}
	c.VideoBackgroundRemovalGreenScreen = &VideoBackgroundRemovalGreenScreenService{client: c}
}
