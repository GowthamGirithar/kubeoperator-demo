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
// the TimeBasedScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *TimeBasedScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	logMap := map[string]interface{}{
		"reqName":           req.Name,
		"reqNamespacedName": req.NamespacedName,
		"reqNamespace":      req.Namespace,
	}
	l.Info(" reconcile info", logMap)

	timeBasedScaleConfig := &scalingv1alpha1.TimeBasedScaler{}
	err := r.Get(ctx, req.NamespacedName, timeBasedScaleConfig)
	if err != nil {
		l.Error(err, "error in getting timeBasedScalerConfig details")
		return ctrl.Result{}, err
	}

	err = r.scaleDeployment(ctx, timeBasedScaleConfig)
	if err != nil {
		timeBasedScaleConfig.Status.Status = "Failure"
		return ctrl.Result{}, err
	}
	timeBasedScaleConfig.Status.Status = "Success"

	// update the status to the custom resource
	err = r.Status().Update(ctx, timeBasedScaleConfig)
	if err != nil {
		l.Error(err, "error in updating the controller status")
	}

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

func (r *TimeBasedScalerReconciler) scaleDeployment(ctx context.Context, timeBasedScaleConfig *scalingv1alpha1.TimeBasedScaler) error {
	l := log.FromContext(ctx)

	currentHour := time.Now().UTC().Hour()
	l.Info("Current Hour", currentHour)
	startHour := timeBasedScaleConfig.Spec.StartHour
	endHour := timeBasedScaleConfig.Spec.EndHour
	if int32(currentHour) > startHour && int32(currentHour) < endHour {
		// get deployment objects for which we need to scale
		deployment := &v1.Deployment{}
		for _, v := range timeBasedScaleConfig.Spec.Deployments {
			err := r.Get(ctx, types.NamespacedName{
				Namespace: v.Namespace,
				Name:      v.Name,
			}, deployment)
			if err != nil {
				l.Error(err, "error in getting deployment details for the deployment ", v.Name)
				return err
			}

			// If we have already configured replicas don't do anything
			// set the labels accordingly
			if deployment.Spec.Replicas == &timeBasedScaleConfig.Spec.ReplicaCount {
				deployment.Labels = map[string]string{
					"time-based-scaling": "false",
				}
				return nil
			} else {
				deployment.Spec.Replicas = &timeBasedScaleConfig.Spec.ReplicaCount
				deployment.Labels = map[string]string{
					"time-based-scaling": "true",
				}
			}

			err = r.Update(ctx, deployment)
			if err != nil {
				l.Error(err, "error in updating deployment details for the deployment ", v.Name)
				return err
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TimeBasedScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1alpha1.TimeBasedScaler{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
