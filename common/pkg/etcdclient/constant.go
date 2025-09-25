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
	keyEtcdProfile = "/cronix/"

	//node节点
	//key /cronix/node/<node_uuid>
	KeyEtcdNodeProfile = keyEtcdProfile + "node/"
	KeyEtcdNode        = KeyEtcdNodeProfile + "%s"

	//key  /cronix/proc/<node_uuid>/<group_id>/<job_id>/<pid>
	KeyEtcdProcProfile = keyEtcdProfile + "proc/%s/"
	KeyEtcdProc        = KeyEtcdProcProfile + "%d/%d/%d"

	//key /cronix/job/<node_uuid>/<group_id>/<job_id>
	KeyEtcdJobProfile = keyEtcdProfile + "job/%s/"
	KeyEtcdJob        = KeyEtcdJobProfile + "%d/%d"

	// key /cronix/once/group/<jobID>
	KeyEtcdOnceProfile = keyEtcdProfile + "once/"
	KeyEtcdOnce        = KeyEtcdOnceProfile + "%d/%d"

	// key /cronix/group/<group_id>
	KeyEtcdGroupProfile = keyEtcdProfile + "group/"
	KeyEtcdGroup        = KeyEtcdGroupProfile + "%d"

	KeyEtcdLock    = keyEtcdProfile + "lock/"
	KeyEtcdNoticer = keyEtcdProfile + "noticer/"
)
