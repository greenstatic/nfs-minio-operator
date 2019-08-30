package nfsminio

import (
	"context"
	"fmt"
	k8v1alpha1 "github.com/greenstatic/nfs-minio-operator/pkg/apis/k8/v1alpha1"
	"golang.org/x/crypto/bcrypt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nfsminio")

// Add creates a new NFSMinio Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNFSMinio{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nfsminio-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NFSMinio
	err = c.Watch(&source.Kind{Type: &k8v1alpha1.NFSMinio{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner NFSMinio
	// Secret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8v1alpha1.NFSMinio{},
	})
	if err != nil {
		return err
	}

	// Pod
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8v1alpha1.NFSMinio{},
	})
	if err != nil {
		return err
	}

	// Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8v1alpha1.NFSMinio{},
	})
	if err != nil {
		return err
	}

	// Ingress
	err = c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8v1alpha1.NFSMinio{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNFSMinio implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNFSMinio{}

// ReconcileNFSMinio reconciles a NFSMinio object
type ReconcileNFSMinio struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NFSMinio object and makes changes based on the state read
// and what is in the NFSMinio.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNFSMinio) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NFSMinio")

	// Fetch the NFSMinio instance
	instance := &k8v1alpha1.NFSMinio{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	namespace := instance.Namespace
	name := instance.Name

	// Check if the secret already exists, if not create a new one
	secret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Secret
		secretKey := randomSecretKey(32)
		sec := r.secretForNFSMinio(instance, instance.Spec.Username, secretKey)
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", sec.Name, "Secret.Name", sec.Name)
		err = r.client.Create(context.TODO(), sec)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Secret", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name)
			return reconcile.Result{}, err
		}

		// Calculate the bcrypt hash of the password, so that if it ever changes we can restart Minio in order to use it.
		bcryptCost := bcrypt.DefaultCost
		hash, err := bcrypt.GenerateFromPassword([]byte(secretKey), bcryptCost)
		if err != nil {
			reqLogger.Error(err, "Failed to use bcrypt hash on new secretKey", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name, "bcryptCost", bcryptCost)
			return reconcile.Result{}, err
		}

		instance.Status.SecretKeyHash = hash
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update NFSMinio status")
			return reconcile.Result{}, err
		}

	} else if err != nil {
		reqLogger.Error(err, "Failed to get Secret")
		return reconcile.Result{}, err
	} else {
		fmt.Println(string(secret.Data["secretKey"]))
	}

	// Check if the deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		dep := r.deploymentForNFSMinio(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Check if the service already exists, if not create a new one
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Service
		serv := r.serviceForNFSMinio(instance)
		reqLogger.Info("Creating a new Service", "Service.Namespace", serv.Namespace, "Service.Name", serv.Name)
		err = r.client.Create(context.TODO(), serv)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", serv.Namespace, "Service.Name", serv.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	// Check if the ingress already exists, if not create a new one
	ingress := &v1beta1.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, ingress)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Ingress
		ing := r.ingressForNFSMinio(instance)
		reqLogger.Info("Creating a new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
		err = r.client.Create(context.TODO(), ing)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Ingress")
		return reconcile.Result{}, err
	}

	//Check if password has changed. If it has kill the Minio pod in order to force using new secrets.
	if err := bcrypt.CompareHashAndPassword(instance.Status.SecretKeyHash, secret.Data["secretKey"]); err != nil {
		reqLogger.Info("Secret key has changed, rehashing password and restarting Minio pod")

		// Kill pod
		podList := &corev1.PodList{}
		labelSelector := labels.SelectorFromSet(labelsForNFSMinio(instance.Name))
		listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
		err = r.client.List(context.TODO(), listOps, podList)
		if err != nil {
			reqLogger.Error(err, "Failed to list pods", "Memcached.Namespace", instance.Namespace, "Memcached.Name", instance.Name)
			return reconcile.Result{}, err
		}

		for _, pod := range podList.Items {
			reqLogger.Info("Killing Minio pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			if err := r.client.Delete(context.TODO(), &pod); err != nil {
				reqLogger.Error(err, "Failed to delete Minio pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			}
		}

		// Recalculate hash and save it
		bcryptCost := bcrypt.DefaultCost
		hash, err := bcrypt.GenerateFromPassword(secret.Data["secretKey"], bcryptCost)
		if err != nil {
			reqLogger.Error(err, "Failed to use bcrypt hash on updated secretKey", "bcryptCost", bcryptCost)
			return reconcile.Result{}, err
		}

		instance.Status.SecretKeyHash = hash
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update NFSMinio status with updated secretKey hash")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// Returns the labels for selecting the resources belonging to the given NFSMinio CR name.
func labelsForNFSMinio(name string) map[string]string {
	return map[string]string{"app": "nfsminio", "nfsminio_cr": name}
}

func (r *ReconcileNFSMinio) secretForNFSMinio(m *k8v1alpha1.NFSMinio, accessKey, secretKey string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			Namespace: m.Namespace,
			Labels: labelsForNFSMinio(m.Name),
		},
		Type: corev1.SecretTypeOpaque,
		StringData:map[string]string{"accessKey": accessKey, "secretKey": secretKey},
	}

	// Set NFSMinio instance as the owner and controller
	controllerutil.SetControllerReference(m, secret, r.scheme)
	return secret
}

func (r *ReconcileNFSMinio) deploymentForNFSMinio(m *k8v1alpha1.NFSMinio) *appsv1.Deployment {
	var replicas int32 = 1
	ls := labelsForNFSMinio(m.Name)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "nfs-data",
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: m.Spec.NFS.Server,
									Path: m.Spec.NFS.Path,
									ReadOnly: m.Spec.NFS.ReadOnly,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: "minio/minio",
						Name: "minio",
						Command: []string{"/bin/sh", "-ce", "/usr/bin/docker-entrypoint.sh minio gateway nas /data"},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name: "nfs-data",
								MountPath: "/data",

							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9000,
						}},
						Env: []corev1.EnvVar{
							{
								Name: "MINIO_ACCESS_KEY",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: m.Name,
										},
										Key: "accessKey",
									},
								},
							},
							{
								Name: "MINIO_SECRET_KEY",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: m.Name,
										},
										Key: "secretKey",
									},
								},
							},
						},
					}},
				},
			},
		},
	}
	// Set NFSMinio instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

func (r *ReconcileNFSMinio) serviceForNFSMinio(m *k8v1alpha1.NFSMinio) *corev1.Service {
	ls := labelsForNFSMinio(m.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			Namespace: m.Namespace,
			Labels: ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port: 80,
					TargetPort: intstr.Parse("9000"),
				},
			},
		},
	}

	// Set NFSMinio instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

func (r *ReconcileNFSMinio) ingressForNFSMinio(m *k8v1alpha1.NFSMinio) *v1beta1.Ingress {
	service := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Name,
			Namespace: m.Namespace,
			Labels: labelsForNFSMinio(m.Name),
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: m.Spec.Domain,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: m.Name,
										ServicePort: intstr.Parse("80"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Set NFSMinio instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
