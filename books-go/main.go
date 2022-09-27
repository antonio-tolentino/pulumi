package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		appName := "books"

		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}

		deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   pulumi.String(appName),
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String("books"),
								Image: pulumi.String("tolentino/books:v1.1.0"),
							}},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		service, err := corev1.NewService(ctx, appName, &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   pulumi.String(appName),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type: pulumi.String("ClusterIP"),
				Ports: &corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(8080),
						Protocol:   pulumi.String("TCP"),
					},
				},
				Selector: appLabels,
			},
		})

		if err != nil {
			return err
		}

		//outputs
		ctx.Export("name", deployment.Metadata.Elem().Name())
		ctx.Export("image", deployment.Spec.Template().Spec().Containers().Index(pulumi.IntInput(pulumi.Int(0))).Image())
		ctx.Export("ip", service.Spec.ClusterIP())

		return nil
	})
}
