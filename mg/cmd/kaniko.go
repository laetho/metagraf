package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/pkg/generators/kaniko"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
)

func init() {
	RootCmd.AddCommand(kanikoCmd)
	kanikoCmd.AddCommand(kanikoBuildCmd)
	kanikoCmd.AddCommand(kanikoCreateCmd)

	kanikoBuildCmd.Flags().StringVarP(&kaniko.KanikoPodOpts.Namespace,"namespace", "n", "", "Provide Kubernets namespace for Pod creation." )
	kanikoBuildCmd.Flags().BoolVar(&Output, "output", false, "Output generated Secret resource.")
	kanikoBuildCmd.Flags().BoolVar(&Dryrun, "dryrun", false, "Settings this to true will not create Secret in kubernetes.")
	kanikoBuildCmd.Flags().BoolVarP(&Watch, "watch", "w", false, "Watch the generated Kaniko Pod.")
	kanikoBuildCmd.Flags().BoolVarP(&Keep, "keep","k",false,"Keep the completed or failed Kaniko Pod." )

	kanikoBuildCmd.Flags().StringVar(&kaniko.KanikoPodOpts.DockerfileArg, "dockerfile", "Dockerfile", "Specify Kaniko --dockerfile argument")
	kanikoBuildCmd.Flags().StringVar(&kaniko.KanikoPodOpts.ContextArg, "context","", "Specify Kaniko --context argument. Overrides Git ref from metaGraf specification.")
	kanikoBuildCmd.Flags().StringVar(&kaniko.KanikoPodOpts.DestinationArg, "destination", "", "Specify Kaniko --destination argument. Registry reference.")

	kanikoCreateCmd.AddCommand(kanikoCreateRegistryCredentialsCmd)
	kanikoCreateRegistryCredentialsCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "Kubernetes namespace for generated Secret.")
	kanikoCreateRegistryCredentialsCmd.Flags().BoolVar(&Output, "output", false, "Output generated Secret resource.")
	kanikoCreateRegistryCredentialsCmd.Flags().BoolVar(&Dryrun, "dryrun", false, "Settings this to true will not create Secret in kubernetes.")
}



var kanikoCmd = &cobra.Command{
	Use:   "kaniko",
	Short: "kaniko operations",
	Long:  MGBanner + ` kaniko `,
}

var kanikoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "kaniko create operations",
	Long:  MGBanner + ` kaniko create`,
}

var kanikoBuildCmd = &cobra.Command{
	Use:   "build <metagraf>",
	Short: "create a kaniko build pod from metaGraf specification",
	Long:  MGBanner + `build kaniko <metagraf.json>`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)

		mg := metagraf.Parse(args[0])

		modules.Variables = GetCmdProperties(mg.GetProperties())
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		generator := kaniko.NewKanikoPodGenerator(mg, metagraf.MGProperties{}, kaniko.KanikoPodOpts)
		obj := generator.GenerateKanikoPod("blah")

		if Output {
			b, err := kaniko.MarshalToYaml(obj)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}

		if !Dryrun {
			err := generator.Create(obj)
			if err != nil {
				log.Fatal(err)
			}
		}

		if Watch {


			stream, err := generator.LogsReader(obj)
			if err != nil {
				log.Fatal(err)
			}
			s := *stream
			defer s.Close()

			for {
				buf := make([]byte, 2000)
				numBytes, err := s.Read(buf)
				if numBytes == 0 {
					continue
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				message := string(buf[:numBytes])
				fmt.Print(message)
			}
		}

		if !Keep {
			err := generator.Delete(obj)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Helper subcommand to create a Kubernetes docker-registry secret from
// interactive input.
var kanikoCreateRegistryCredentialsCmd = &cobra.Command{
	Use:   "registry-credentials <secretname>",
	Short: "creates a docker-registry secret ",
	Long:  "Creates a Kubernetes docker-registry secret interactively",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			log.Fatal(errNoSecretNameProvided)
		}

		var (
			regsrv string
			reguser string
			regpass string
			regemail string
			// BASE64 encoding of "username:password"
			regauth string
		)

		fmt.Println("Hi there! I'll help you create a Kubernetes Secret of the docker-registry type.")
		fmt.Println("Enter registry server: ")
		regsrv = readInput()
		fmt.Println("Enter registry user: ")
		reguser = readInput()
		fmt.Println("Enter registry password or token: ")
		regpass = readInput()
		fmt.Println("Enter email associated with registry if applicable: ")
		regemail = readInput()

		regauth = base64.StdEncoding.EncodeToString([]byte(reguser+":"+regpass))

		// server, user, password, email, base64(auth(username:password))
		jsonstring := fmt.Sprintf("{\"auths\":{\"%v\":{\"username\":\"%v\",\"password\":\"%v\",\"email\":\"%v\",\"auth\":\"%v\"}}}", regsrv,reguser,regpass,regemail,regauth)

		stringdata := make(map[string]string)
		stringdata[corev1.DockerConfigJsonKey] = jsonstring

		sec := corev1.Secret{
			TypeMeta:   metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			StringData: stringdata,
			Type:       corev1.SecretTypeDockerConfigJson,
		}
		sec.ObjectMeta.Name = args[0]
		sec.ObjectMeta.Namespace = Namespace

		if !Dryrun {
			client := k8sclient.GetCoreClient()
			res, err := client.Secrets(Namespace).Create(context.TODO(), &sec, metav1.CreateOptions{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("created secret: %v", res.Name)
		}

		if Output {
			o, _ := yaml.Marshal(sec)
			fmt.Println(string(o))
		}
	},
}

// Reads input from stdin and returns string after newline(enter).
func readInput() string {
	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n\r")
	return text
}