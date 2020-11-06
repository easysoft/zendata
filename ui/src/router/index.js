import Vue from "vue";
import VueRouter from "vue-router";

import MainLayout from "../layout/MainLayout";
import ColLayout from "../layout/ColLayout";

Vue.use(VueRouter);

const routes = [
    {
        path:"/",
        name: "",
        component: MainLayout,
        redirect: '/data/mine',
        children: [
            {
                path:"data",
                name: "data",
                component: ColLayout,
                children: [
                    {
                        path: 'mine',
                        name: 'mine',
                        component: () => import('../views/data/Mine')
                    },
                    {
                        path: 'excel',
                        name: 'excel',
                        component: () => import('../views/data/Excel')
                    },
                    {
                        path: 'yaml',
                        name: 'yaml',
                        component: () => import('../views/data/Yaml')
                    },
                ]
            }
        ]
    },
]

const router =  new VueRouter({
    routes
})
export default router;
