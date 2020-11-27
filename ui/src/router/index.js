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
        redirect: '/data/mine/list',
        children: [
            {
                path:"data",
                name: "data",
                component: ColLayout,
                children: [
                    {
                        path: 'mine',
                        name: 'mine',
                        component: BlankLayout,
                        children: [
                            {
                                path: 'list',
                                name: 'mine-list',
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
                        children: [
                            {
                                path: 'ranges',
                                name: 'ranges',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list',
                                        name: 'ranges-list',
                                        component: () => import('../views/data/buildin/ranges/List')
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'ranges-edit',
                                        component: () => import('../views/data/buildin/ranges/Edit')
                                    },
                                ],
                            },
                            {
                                path: 'instances',
                                name: 'instances',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list',
                                        name: 'instances-list',
                                        component: () => import('../views/data/buildin/instances/List')
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'instances-edit',
                                        component: () => import('../views/data/buildin/instances/Edit')
                                    },
                                ],
                            },
                            {
                                path: 'excel',
                                name: 'excel',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list',
                                        alias: "index",
                                        name: 'excel-list',
                                        component: () => import('../views/data/buildin/excel/List')
                                    },
                                    {
                                        path: 'excel/:id',
                                        name: 'excel-edit',
                                        component: () => import('../views/data/buildin/excel/Edit')
                                    },
                                ],
                            },
                            {
                                path: 'config',
                                name: 'config',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list',
                                        name: 'config-list',
                                        component: () => import('../views/data/buildin/config/List')
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'config-edit',
                                        component: () => import('../views/data/buildin/config/Edit')
                                    },
                                ],
                            },
                            {
                                path: 'text',
                                name: 'text',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list',
                                        name: 'text-list',
                                        component: () => import('../views/data/buildin/text/List')
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'text-edit',
                                        component: () => import('../views/data/buildin/text/Edit')
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
    mode: 'history',
    routes
})
export default router;
