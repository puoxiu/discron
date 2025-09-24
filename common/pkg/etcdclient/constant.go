package etcdclient

/*
   "node": "/cronix/node/",
   "proc": "/cronix/proc/",
   "cmd": "/cronix/cmd/",
   "once": "/cronix/once/",
   "lock": "/cronix/lock/",
   "group": "/cronix/group/",
   "noticer": "/cronix/noticer/"
*/
const (
	KeyEtcdProfile = "/cronix/"
	//node节点
	KeyEtcdNode    = KeyEtcdProfile + "node/"
	KeyEtcdProc    = KeyEtcdProfile + "proc/"
	KeyEtcdCmd     = KeyEtcdProfile + "cmd/"
	KeyEtcdOnce    = KeyEtcdProfile + "once/"
	KeyEtcdLock    = KeyEtcdProfile + "lock/"
	KeyEtcdNoticer = KeyEtcdProfile + "noticer/"
)
