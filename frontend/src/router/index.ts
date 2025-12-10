import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Layouts
import MainLayout from '@/layouts/MainLayout.vue';
import DashboardLayout from '@/layouts/DashboardLayout.vue';
import AdminLayout from '@/layouts/AdminLayout.vue';

// Public pages
import HomePage from '@/views/HomePage.vue';
import LoginPage from '@/views/auth/LoginPage.vue';
import RegisterPage from '@/views/auth/RegisterPage.vue';
import ActivatePage from '@/views/auth/ActivatePage.vue';

// Dashboard pages
import DashboardPage from '@/views/dashboard/DashboardPage.vue';
import LinkGroupsPage from '@/views/dashboard/LinkGroupsPage.vue';
import LinkGroupFormPage from '@/views/dashboard/LinkGroupFormPage.vue';

// Admin pages
import AdminDashboard from '@/views/admin/AdminDashboard.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Main layout (для неавторизованных пользователей)
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: false },
      children: [
        {
          path: '',
          name: 'home',
          component: HomePage
        },
        {
          path: 'auth/login',
          name: 'login',
          component: LoginPage
        },
        {
          path: 'auth/register',
          name: 'register',
          component: RegisterPage
        },
        {
          path: 'auth/confirm/:token',
          name: 'activate',
          component: ActivatePage,
          props: true
        }
      ]
    },

    // Dashboard layout (для авторизованных пользователей)
    {
      path: '/dashboard',
      component: DashboardLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: DashboardPage
        },
        {
          path: 'link-groups',
          name: 'link-groups',
          component: LinkGroupsPage
        },
        {
          path: 'link-groups/create',
          name: 'link-group-create',
          component: LinkGroupFormPage
        },
        {
          path: 'link-groups/edit/:id',
          name: 'link-group-edit',
          component: LinkGroupFormPage,
          props: true
        }
      ]
    },

    // Admin layout (только для админов)
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: AdminDashboard
        }
      ]
    },

    // Redirects
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
});

// Guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();

  // Инициализируем store если нужно
  if (!authStore.user && localStorage.getItem('access_token')) {
    await authStore.loadProfile();
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin);

  // Проверка авторизации
  if (requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login');
    return;
  }

  // Проверка прав админа
  if (requiresAdmin && !authStore.isAdmin) {
    next('/dashboard');
    return;
  }

  // Если пользователь авторизован и пытается зайти на публичные страницы
  if (authStore.isAuthenticated &&
    (to.name === 'login' || to.name === 'register' || to.name === 'home')) {
    next('/dashboard');
    return;
  }

  next();
});

export default router;
