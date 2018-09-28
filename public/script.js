const App = async() => {
    const template = await axios.get('/app.html').then(res => {
        return res.data;
    });

    return {
        template,

        created() {
            this.checkCurrentLogin()
        },

        updated() {
            this.checkCurrentLogin()
        },

        methods: {
            checkCurrentLogin() {
                if (!this.$app.currentUser && this.$route.path !== '/login') {
                    this.$router.push('/login?redirect=' + this.$route.path)
                }
            }
        }
    };
}

const Login = async() => {
    const template = await axios.get('/login.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        data() {
            return {
                username: '',
                password: ''
            }
        },
        
        methods: {
            async login(e) {
                let form = { username: this.username, password: this.password }
                axios.post('/login', form).then(res => {
                    this.$app.currentUser = res.data
                    this.$router.push(this.$route.query.redirect)
                }).catch((e) => {
                    console.log(e)
                });

                e.preventDefault();
            }
        }
    }
}

const Overview = async() => {
    const template = await axios.get('/overview.html').then(res => {
        return res.data;
    });
    
    return {
        template,
        
        data() {
            return {
            };
        },

        async created() {
        },
        
        methods: {
        }
    };
}

const DeviceList = async () => {
    const template = await axios.get('/device_list.html').then(res => {
        return res.data;
    });
    
    return {
        template,
        
        data() {
            return {
                devices: []
            };
        },

        async created() {
            this.devices = await this.getDevices();
        },
        
        methods: {
            async getDevices() {
                return await axios.get("/api/devices").then(res => {
                    return res.data;
                });
            }
        },

        filters: {
            pretty(v) {
                return v.substring(0, 6) + "..."
            }
        }
    };
}

const DeviceView = async () => {
    const template = await axios.get('/device_view.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        props: ['uid', 'embedded'],
        
        data() {
            return {
                device: {}
            }
        },

        async created() {
            let uid = this.$route.params.uid || this.uid
            this.device = await this.getDevice(uid);
        },
        
        methods: {
            async getDevice(uid) {
                return await axios.get("/api/devices/" + uid).then(res => {
                    return res.data;
                });
            }
        },

        filters: {
            pretty: function(value) {
                return JSON.stringify(value, null, 2).replace(/\n/g, "<br/>").replace(/ /g,"&nbsp;")
            }
        }
    }
}

const PackageList = async () => {
    const template = await axios.get('/package_list.html').then(res => {
        return res.data;
    });
    
    return {
        template,
        
        data() {
            return {
                packages: []
            };
        },

        async created() {
            this.packages = await this.getPackages();
        },
        
        methods: {
            async getPackages() {
                return await axios.get("/api/packages").then(res => {
                    return res.data;
                });
            },

            uploadPackage(e) {
                let form = new FormData();
                form.append('file', this.$refs.file.files[0]);

                axios.post('/api/packages', form,
                           { headers: { 'Content-Type': 'multipart/form-data' } }).then(function(){
                               console.log('SUCCESS!!');
                           })
                    .catch(function(){
                        console.log('FAILURE!!');
                    });
            } 
        }
    };
}

const PackageView = async () => {
    const template = await axios.get('/package_view.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        data() {
            return {
                pkg: { supported_hardware: [], metadata: { objects: [[]] } },
                selectedObject: -1,
                currentInstallSet: 0,
            }
        },

        async created() {
            this.pkg = await this.getPackage();
            this.pkg.metadata = JSON.parse(window.atob(this.pkg.metadata))
        },
        
        methods: {
            async getPackage() {
                return await axios.get("/api/packages/" + this.$route.params.uid).then(res => {
                    return res.data;
                });
            }
        },

        filters: {
            pretty: function(value) {
                return JSON.stringify(value, null, 2).replace(/\n/g, "<br/>").replace(/ /g,"&nbsp;")
            }
        }
    }
}

const RolloutList = async () => {
    const template = await axios.get('/rollout_list.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        data() {
            return {
                rollouts: [],
            }
        },

        async created() {
            this.rollouts = await this.getRollouts();
            this.rollouts.forEach(async rollout => {
                rollout.package = await this.getPackage(rollout.package)
            })
        },
        
        methods: {
            async getRollouts() {
                return await axios.get("/api/rollouts").then(res => {
                    return res.data;
                });
            },

            async getPackage(uid) {
                return await axios.get("/api/packages/" + uid).then(res => {
                    return res.data;
                });
            }
        },

        filters: {
            humanizeDate(v) {
                return moment(v).format('lll')
            }
        }
    }
}

const RolloutView = async () => {
    const template = await axios.get('/rollout_view.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        components: { 'Chato': Chato },
        
        data() {
            return {
                rollout: { package: {}, devices: [], statistics: { }, status: '' },
                timer: null,
            }
        },
        
        async created() {
            this.refresh()

            this.timer = setInterval(this.refresh, 5 * 1000)
        },

        beforeDestroy() {
            clearInterval(this.timer)
        },
        
        methods: {
            async refresh() {
                this.rollout = await this.getRollout().then(async rollout => {
                    rollout.package = await this.getPackage(rollout.package)
                    rollout.statistics = await this.getStatistics(rollout)

                    if (rollout.finished_at > rollout.started_at) {
                        rollout.status = "finished"
                    } else if (rollout.running) {
                        rollout.status = "running"
                    } else {
                        rollout.status = "paused"
                    }

                    return rollout
                })
            },
            
            async getRollout() {
                return await axios.get("/api/rollouts/" + this.$route.params.id).then(res => {
                    return res.data;
                });
            },

            async getPackage(uid) {
                return await axios.get("/api/packages/" + uid).then(res => {
                    return res.data;
                });
            },

            async getStatistics(rollout) {
                return await axios.get("/api/rollouts/" + rollout.id + "/statistics").then(res => {
                    return res.data;
                });
            }
        },

        filters: {
            pretty: function(value) {
                return JSON.stringify(value, null, 2).replace(/\n/g, "<br/>").replace(/ /g,"&nbsp;")
            },

            humanizedDate(v) {
                return moment(v).format('lll')
            }
        }
    }
}

const Chato = async() => {
    return {
        extends: VueChartJs.Doughnut,

        mounted() {
            this.renderChart({
                height: 200,
                labels: ['Pending', 'In Progress', 'Finished', 'Failed'],
                datasets: [
                    {
                        data: [25, 25, 25, 25],
                        backgroundColor: ['#ddd','yellow','#58ED93','#F94D35']
                    }
                ]
            }, {responsive: true, maintainAspectRatio: true, legend: { display: false }, cutoutPercentage: 80})
        }
    }
}

const RolloutNew = async () => {
    const template = await axios.get('/rollout_new.html').then(res => {
        return res.data;
    });
    
    return {
        template,

        components: { 'DeviceView': DeviceView },

        data() {
            return {
                selectedPackage: null,
                packages: [],
                devices: [],
                opened: "",
                selectedDevices: [],
            }
        },

        async created() {
            this.refresh()
        },
        
        methods: {
            async refresh() {
                this.devices = await this.getDevices()
                this.packages = await this.getPackages()
            },
        
            async getDevices() {
                return await axios.get("/api/devices").then(res => {
                    return res.data;
                });
            },

            async getPackages() {
                return await axios.get("/api/packages").then(res => {
                    return res.data;
                });
            },

            toggleOpened(uid) {
                if (this.opened == uid) {
                    this.opened = "";
                } else {
                    this.opened = ""
                    this.$nextTick(() => {
                        this.opened = uid
                    })
                }
            },

            toggleDeviceSelection(uid) {
                let index = this.selectedDevices.findIndex((device, i) => {
                    return device.uid == uid
                })

                if (index >= 0) {
                    this.selectedDevices.splice(index, 1)
                } else {
                    this.selectedDevices.push(this.compatibleDevices().find(d => { return d.uid == uid }))
                }
            },

            selectAllDevices() {
                this.selectedDevices = this.compatibleDevices()
            },

            compatibleDevices() {
                return this.devices.filter(device => {
                    return device.version != this.selectedPackage.version && (this.selectedPackage.supported_hardware == "any" || this.selectedPackage.supported_hardware.includes(device.hardware)) })
            },
            
            save(start) {
                axios.post('/api/rollouts', {
                    'package': this.selectedPackage.uid,
                    devices: this.selectedDevices.map((d) => { return d.uid }),
                    running: start }).then(res => {
                        let rollout = res.data;
                        this.$router.push('/rollouts/' + rollout.id)
                    }).catch((e) => {
                        console.log(e)
                    });
            },
        },

        filters: {
            differentFromPackage(devices, selectedPackage) {
                return devices.filter(device => {
                    return device.version != selectedPackage.version && (selectedPackage.supported_hardware == "any" || selectedPackage.supported_hardware.includes(device.hardware)) })
            }
        }
    }
}

const routes = [
    { path: "/", redirect: '/overview' },
    { path: "/login", component: Login },
    { path: "/overview", component: Overview },
    { path: "/devices", component: DeviceList },
    { path: "/devices/:uid", component: DeviceView },
    { path: "/packages", component: PackageList },
    { path: "/packages/:uid", component: PackageView },
    { path: "/rollouts", component: RolloutList },
    { path: "/rollouts/new", component: RolloutNew },
    { path: "/rollouts/:id", component: RolloutView },
];

const router = new VueRouter({
    routes
});

window.onload = function() {
    var app = new Vue({
        el: '#app',
        router,
        template: '<App/>',
        components: { App, Login },
        computed: {
            currentUser: {
                cache: false,
                
                get() {
                    if (!localStorage.currentUser) {
                        return null
                    } else {
                        return JSON.parse(localStorage.currentUser)
                    }
                },

                set(currentUser) {
                    localStorage.currentUser = JSON.stringify(currentUser)
                    axios.defaults.headers.common['Authorization'] =  'Bearer ' + currentUser.token
                }
            }
        }
    });

    if (app.currentUser) {
        axios.defaults.headers.common['Authorization'] =  'Bearer ' + app.currentUser.token
    }
    

    Vue.prototype.$app = app
};
