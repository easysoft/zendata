import Vue from "vue";
import VueRouter from "vue-router";

import MainLayout from "../layout/MainLayout";
import ColLayout from "../layout/ColLayout";
import BlankLayout from "../layout/BlankLayout";
import BuildinLayout from "../views/data/buildin/Layout";

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
                                component: () => import('../views/data/mine/List'),
                                meta: { title: 'menu.data.list' }
                            },
                            {
                                path: 'edit/:id',
                                name: 'mine-edit',
                                component: () => import('../views/data/mine/Edit'),
                                meta: { title: 'menu.data.edit' }
                            },
                        ],
                    },
                    {
                        path: 'buildin',
                        name: 'buildin',
                        component: BuildinLayout,
                        children: [
                            {
                                path: 'ranges',
                                name: 'ranges',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list/:id?',
                                        name: 'ranges-list',
                                        component: () => import('../views/data/buildin/ranges/List'),
                                        meta: { title: 'menu.ranges.list' }
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'ranges-edit',
                                        component: () => import('../views/data/buildin/ranges/Edit'),
                                        meta: { title: 'menu.ranges.edit' }
                                    },
                                ],
                            },
                            {
                                path: 'instances',
                                name: 'instances',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list/:id?',
                                        name: 'instances-list',
                                        component: () => import('../views/data/buildin/instances/List'),
                                        meta: { title: 'menu.instances.list' }
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'instances-edit',
                                        component: () => import('../views/data/buildin/instances/Edit'),
                                        meta: { title: 'menu.instances.edit' }
                                    },
                                ],
                            },
                            {
                                path: 'excel',
                                name: 'excel',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list/:id?',
                                        name: 'excel-list',
                                        component: () => import('../views/data/buildin/excel/List'),
                                        meta: { title: 'menu.excel.list' }
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'excel-edit',
                                        component: () => import('../views/data/buildin/excel/Edit'),
                                        meta: { title: 'menu.excel.edit' }
                                    },
                                ],
                            },
                            {
                                path: 'config',
                                name: 'config',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list/:id?',
                                        name: 'config-list',
                                        component: () => import('../views/data/buildin/config/List'),
                                        meta: { title: 'menu.config.list' }
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'config-edit',
                                        component: () => import('../views/data/buildin/config/Edit'),
                                        meta: { title: 'menu.config.edit' }
                                    },
                                ],
                            },
                            {
                                path: 'text',
                                name: 'text',
                                component: BlankLayout,
                                children: [
                                    {
                                        path: 'list/:id?',
                                        name: 'text-list',
                                        component: () => import('../views/data/buildin/text/List'),
                                        meta: { title: 'menu.text.list' }
                                    },
                                    {
                                        path: 'edit/:id',
                                        name: 'text-edit',
                                        component: () => import('../views/data/buildin/text/Edit'),
                                        meta: { title: 'menu.text.edit' }
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
    mode: 'hash',
    routes
})
export default router;
