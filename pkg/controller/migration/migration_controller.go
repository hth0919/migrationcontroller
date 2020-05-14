package migration

import (
	"context"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	logs "log"
	"os"
	"time"

	ketiv1alpha1 "github.com/hth0919/migrationcontroller/pkg/apis/keti/v1alpha1"

	cp "github.com/hth0919/checkpointproto"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"gopkg.in/yaml.v2"
)

var log = logf.Log.WithName("controller_migration")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Migration Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMigration{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("migration-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Migration
	err = c.Watch(&source.Kind{Type: &ketiv1alpha1.Migration{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Migration
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ketiv1alpha1.Migration{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMigration implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMigration{}

// ReconcileMigration reconciles a Migration object
type ReconcileMigration struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Migration object and makes changes based on the state read
// and what is in the Migration.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMigration) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Migration")

	// Fetch the Migration instance
	migration := &ketiv1alpha1.Migration{}
	err := r.client.Get(context.TODO(), request.NamespacedName, migration)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if migration.Spec.Purpose == "Convert" || migration.Spec.Purpose == "convert" {
		migration.Spec.Purpose = "convert"
		result, err := convertHandler(migration,r)
		if err != nil {
			reqLogger.Error(err, err.Error())
			return result, err
		}
	} else if migration.Spec.Purpose == "Checkpoint" || migration.Spec.Purpose == "checkpoint" || migration.Spec.Purpose == "CheckPoint" {
		migration.Spec.Purpose = "checkpoint"
		result, err := checkpointHandler(migration,r)
		if err != nil {
			reqLogger.Error(err, err.Error())
			return result, err
		}
	} else if migration.Spec.Purpose == "Migration" || migration.Spec.Purpose == "migration" {
		migration.Spec.Purpose = "migration"
		result, err := migrationHandler(migration,r)
		if err != nil {
			reqLogger.Error(err, err.Error())
			return result, err
		}
	} else {
		reqLogger.Error(err, "Unsupported Value", migration.Spec.Purpose, " in keti.Migration")
		return reconcile.Result{}, err
	}
/*

		found := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: migration.Name, Namespace: migration.Namespace}, found)
	if err != nil && errors.IsNotFound(err) { // 만약에 migration 위한 Deployment가 없다면
		// 새로운 Deployment를 생성합니다. deploymentForMemcached 함수는 Deployment를 위한 spec을 반환합니다.
		dem := r.daemonSetForMigration(migration)
		reqLogger.Info("Creating a new DaemonSet", "DaemonSet.Namespace", dem.Namespace, "DaemonSet.Name", dem.Name)
		err = r.client.Create(context.TODO(), dem)
		if err != nil {
			//reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dem.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment가 성공적으로 생성되었다면, 이 이벤트를 다시 Requeue 합니다.
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	} else if err == nil {

	}*/

	return reconcile.Result{}, nil
}

func convertHandler(m *ketiv1alpha1.Migration, r *ReconcileMigration) (reconcile.Result, error) {
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: m.Spec.Podname, Namespace: m.Spec.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	m.Spec.Pod.Type = found.TypeMeta
	m.Spec.Pod.Object = found.ObjectMeta
	m.Spec.Pod.PodSpec = found.Spec

	/*체크포인트 생성 코드*/
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	node := &corev1.Node{}
	nodeIP := ""
	node, err = clientset.CoreV1().Nodes().Get(m.Spec.DestinationNode,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	for i := 0;i<len(node.Status.Addresses);i++ {
		if node.Status.Addresses[i].Type == "InternalIP" {
			nodeIP = node.Status.Addresses[i].Address
		}
	}
	Host := nodeIP + ":10350"
	conn, err := grpc.Dial(Host, grpc.WithInsecure())
	if err != nil {
		logs.Fatalln("did not connect: ", err)
	}

	defer conn.Close()
	c := cp.NewCheckpointPeriodClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	createctruct := &cp.CreateCheckpoint{
		PodName:              &found.Name,
	}
	_ , err = c.CheckpointCreate(ctx, createctruct)
	if err != nil {
		logs.Fatalln("did not connect: ", err)
	}

	m.Spec.DestinationNode = m.Spec.Node
	return migrationHandler(m,r)
}

func checkpointHandler(m *ketiv1alpha1.Migration, r *ReconcileMigration) (reconcile.Result, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	node := &corev1.Node{}
	nodeIP := ""
	node, err = clientset.CoreV1().Nodes().Get(m.Spec.DestinationNode,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	for i := 0;i<len(node.Status.Addresses);i++ {
		if node.Status.Addresses[i].Type == "InternalIP" {
			nodeIP = node.Status.Addresses[i].Address
		}
	}
	Host := nodeIP + ":10350"
	conn, err := grpc.Dial(Host, grpc.WithInsecure())
	if err != nil {
		logs.Fatalln("did not connect: ", err)
	}

	defer conn.Close()
	c := cp.NewCheckpointPeriodClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &cp.InputValue{
		Period:               &m.Spec.Period,
		PodName:              nil,
	}
	_ , err = c.SetCheckpointPeriod(ctx, in)
	if err != nil {
		logs.Fatalln("did not connect: ", err)
	}
	return reconcile.Result{Requeue: true}, nil
}

func migrationHandler(m *ketiv1alpha1.Migration, r *ReconcileMigration) (reconcile.Result, error) {
	pod := &corev1.Pod{
		TypeMeta:   m.Spec.Pod.Type,
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       m.Spec.Pod.PodSpec,
		Status:     corev1.PodStatus{},
	}

	y, err := yaml.Marshal(pod)
	if err != nil {
		log.Error(err, "err")
	}
	dir := "/nfs/"+m.Spec.DestinationNode+"/"+m.Spec.Podname+".yaml"
	origindir := "/nfs/"+m.Spec.Node+"/"+m.Spec.Podname+".yaml"
	//파일 쓰기
	err = ioutil.WriteFile(dir, y, 0)
	if err != nil {
		panic(err)
	}
	if m.Spec.Node!=m.Spec.DestinationNode {
		err3 := os.Remove(origindir)
		if err3 != nil {
			panic(err3)
		}
	}else {
		if m.Spec.Purpose == "convert" {
			return podDeleteHandler(m,r)
		}
	}

	return reconcile.Result{Requeue: true}, nil
}

func podDeleteHandler(m *ketiv1alpha1.Migration, r *ReconcileMigration) (reconcile.Result, error) {
	found := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: m.Spec.Podname, Namespace: m.Spec.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{Requeue: true}, nil
}
/*

func labelsForMigration(name string) map[string]string {
	return map[string]string{"app": "migration", "migration_cr": name}
}*/