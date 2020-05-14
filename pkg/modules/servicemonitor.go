package modules

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/pkg/metagraf"
)

func GenServiceMonitor(mg *metagraf.MetaGraf) {
	sm := monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.ObjectMeta.Name,
			Namespace: s.ObjectMeta.Namespace,
			Labels:    labels,
		},
	}
	print(sm)
/*
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.ObjectMeta.Name,
			Namespace: s.ObjectMeta.Namespace,
			Labels:    labels,
			},
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
			Endpoints: endpoints,
		},
	}
*/
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