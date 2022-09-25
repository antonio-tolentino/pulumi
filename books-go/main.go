package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		appName := pulumi.String("books")
		appLabels := pulumi.StringMap{
			"app": pulumi.String("books"),
		}
		deployment, err := appsv1.NewDeployment(ctx, "books", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   appName,
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

		ctx.Export("name", deployment.Metadata.Elem().Name())
		ctx.Export("image", deployment.Spec.Template().Spec().Containers().Index(pulumi.IntInput(pulumi.Int(0))).Image())

		return nil
	})
}
