package modules

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/mg/k8sclient"
	"metagraf/mg/params"
	"metagraf/pkg/metagraf"
	"fmt"
	"os"
	log "k8s.io/klog"
)

func GenServiceMonitorAndService(mg *metagraf.MetaGraf) {

}

func GenServiceMonitor(mg *metagraf.MetaGraf, svc *corev1.Service) {
	objname := Name(mg)
	// Resource labels
	l := make(map[string]string)
	l["app"] = objname
	l["deploymentconfig"] = objname
	l["prometheus"] = params.ServiceMonitorOperatorName

	// Selector
	s := make(map[string]string)
	s["app"] = objname

	var metricPort corev1.ServicePort
	for _,p := range svc.Spec.Ports {
		if p.Port == params.ServiceMonitorPort {
			metricPort = p
		}
	}
	if metricPort.Port != params.ServiceMonitorPort {
		fmt.Printf("ERROR: Unable to find default or provided service port for scraping. Tried to find: %v", params.ServiceMonitorPort)
		os.Exit(1)
	}

	eps := []monitoringv1.Endpoint{}
	ep := monitoringv1.Endpoint{
		Port:                 metricPort.Name,
		Path:                 params.ServiceMonitorPath,
		Scheme:               params.ServiceMonitorScheme,
		Interval:             params.ServiceMonitorInterval,
	}
	eps = append(eps, ep)

	var obj = monitoringv1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceMonitor",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      objname,
			Namespace: NameSpace,
			Labels:    l,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: eps,
			Selector: metav1.LabelSelector{
				MatchLabels: l,
			},
		},
	}

	if !Dryrun {
		StoreServiceMonitor(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}

}

func StoreServiceMonitor(obj monitoringv1.ServiceMonitor) {
	client := k8sclient.GetMonitoringV1Client().ServiceMonitors(NameSpace)
	res, _ := client.Get(obj.Name, metav1.GetOptions{})

	if len(res.ResourceVersion) > 0 {
		obj.ResourceVersion = res.ResourceVersion
		_, err := client.Update(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated ServiceMonitor: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created ServiceMonitor: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteServiceMonitor(name string) {
	client := k8sclient.GetMonitoringV1Client().ServiceMonitors(NameSpace)

	_, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("The ServiceMonitor: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete ServiceMonitor: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted ServiceMonitor: ", name, ", in namespace: ", NameSpace)
}