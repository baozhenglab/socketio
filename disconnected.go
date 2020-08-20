package sckio

func (hdl *socketHandler) onDisconnected(as *appSocket) {
	if as.cu != nil {
		hdl.UserDisconnect(as.cu.UserID())
	}
	as.Logger().Debugf("socket %s disconnected", as)
}
