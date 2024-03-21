/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	scalingv1alpha1 "github.com/GowthamGirithar/kubeoperator-demo/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TimeBasedScalerReconciler reconciles a TimeBasedScaler object
type TimeBasedScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=scaling.gg.com,resources=timebasedscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scaling.gg.com,resources=timebasedscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scaling.gg.com,resources=timebasedscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TimeBasedScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *TimeBasedScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info(" inside reconcile loop ", "req name :", req.Name, ",req.NamespacedName :", req.NamespacedName, ", req.Namespace", req.Namespace)
	deploymentConfig := &scalingv1alpha1.TimeBasedScaler{}
	err := r.Get(ctx, req.NamespacedName, deploymentConfig)
	if err != nil {
		l.Error(err, "error in getting crd details")
		return ctrl.Result{}, err
	}
	currentHr := time.Now().UTC().Hour()
	if int32(currentHr) > deploymentConfig.Spec.StartHour && int32(currentHr) < deploymentConfig.Spec.EndHour {
		deployment := &v1.Deployment{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: deploymentConfig.Spec.Deployments[0].Namespace,
			Name:      deploymentConfig.Spec.Deployments[0].Name,
		}, deployment)
		if err != nil {
			l.Error(err, "error in getting deployment details")
			return ctrl.Result{}, err
		}

		if deployment.Spec.Replicas == &deploymentConfig.Spec.ReplicaCount {
			return ctrl.Result{}, nil
		}

		deployment.Spec.Replicas = &deploymentConfig.Spec.ReplicaCount
		err = r.Update(ctx, deployment)
		if err != nil {
			deploymentConfig.Status.Status = "Failed"
			l.Error(err, "error in updating deployment details")
			return ctrl.Result{}, err
		}
	}
	deploymentConfig.Status.Status = "Success"
	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TimeBasedScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1alpha1.TimeBasedScaler{}).
		Complete(r)
}
