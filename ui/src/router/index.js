import Vue from "vue";
import VueRouter from "vue-router";

import MainLayout from "@/layout/MainLayout";

Vue.use(VueRouter);

const routes = [
    {
        path:"/",
        name: "",
        component: MainLayout,
        redirect: '/index',
        children: [
            {
                path: '/index',
                name: 'index',
                component: () => import('@/views/Index')
            },
            {
                path: '/test',
                name: 'test',
                component: () => import('@/views/test/Test')
            }
        ]
    }
]

const router =  new VueRouter({
    routes
})
export default router;
