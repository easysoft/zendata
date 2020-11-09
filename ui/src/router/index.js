import Vue from "vue";
import VueRouter from "vue-router";

import MainLayout from "../layout/MainLayout";
import ColLayout from "../layout/ColLayout";
import BlankLayout from "../layout/BlankLayout";

Vue.use(VueRouter);

const routes = [
    {
        path:"/",
        name: "",
        component: MainLayout,
        redirect: '/data/mine/index',
        children: [
            {
                path:"data",
                name: "data",
                component: ColLayout,
                redirect: '/data/mine/index',
                children: [
                    {
                        path: 'mine',
                        name: 'mine',
                        component: BlankLayout,
                        redirect: '/data/mine/index',
                        children: [
                            {
                                path: 'mine-list',
                                alias: "index",
                                name: 'list',
                                component: () => import('../views/data/mine/List')
                            },
                            {
                                path: 'edit/:id',
                                name: 'mine-edit',
                                component: () => import('../views/data/mine/Edit')
                            },
                        ],
                    },
                    {
                        path: 'buildin',
                        name: 'buildin',
                        component: BlankLayout,
                        redirect: '/data/buildin/excel/index',
                        children: [
                            {
                                path: 'excel',
                                name: 'excel',
                                component: BlankLayout,
                                redirect: '/data/buildin/excel/index',
                                children: [
                                    {
                                        path: 'list',
                                        alias: "index",
                                        name: 'excel-list',
                                        component: () => import('../views/data/buildin/excel/List')
                                    },
                                ],
                            },
                            {
                                path: 'yaml',
                                name: 'yaml',
                                component: BlankLayout,
                                redirect: '/data/buildin/yaml/index',
                                children: [
                                    {
                                        path: 'list',
                                        alias: "index",
                                        name: 'yaml-list',
                                        component: () => import('../views/data/buildin/yaml/List')
                                    },
                                ],
                            },
                        ],
                    }
                ]
            }
        ]
    },
]

const router =  new VueRouter({
    routes
})
export default router;
