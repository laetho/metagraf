package modules

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/pkg/metagraf"
)

func GenServiceMonitor(mg *metagraf.MetaGraf) {
	objname := Name(mg)
	// Resource labels
	l := make(map[string]string)
	l["app"] = objname
	l["deploymentconfig"] = objname

	// Selector
	s := make(map[string]string)
	s["app"] = objname
	s["deploymentconfig"] = objname


	// todo: we need to detect this from service object
	eps := []monitoringv1.Endpoint{}
	ep := monitoringv1.Endpoint{
		Port:                 "8080",
		Path:                 "/metrics",
		Scheme:               "http",
		Interval:             "30s",
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
/*	client := ocpclient.GetBuildClient().BuildConfigs(NameSpace)
	bc, _ := client.Get(obj.Name, metav1.GetOptions{})

	if len(bc.ResourceVersion) > 0 {
		obj.ResourceVersion = bc.ResourceVersion
		_, err := client.Update(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated BuildConfig: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created BuildConfig: ", obj.Name, " in Namespace: ", NameSpace)
	}
 */
}

func DeleteServiceMonitor(name string) {
	/*
	client := ocpclient.GetBuildClient().BuildConfigs(NameSpace)

	_, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("The BuildConfig: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete BuildConfig: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted BuildConfig: ", name, ", in namespace: ", NameSpace)
	 */
}