import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * The curriculum reads top-to-bottom: Start Here → Roadmap → Modules 0–4 →
 * Capstone. Each module is a sidebar category with a generated landing page,
 * and every page inside a module is numbered so the sidebar order matches the
 * intended study order.
 */
const sidebars: SidebarsConfig = {
  curriculumSidebar: [
    'index',
    'roadmap',
    {
      type: 'category',
      label: 'Module 0 — Setup & Orientation',
      collapsed: false,
      link: {
        type: 'generated-index',
        title: 'Module 0 — Setup & Orientation',
        description:
          'Week 1. Install the toolchain, run the full Data Exchange stack locally, and build a mental map of the platform you are about to learn.',
        slug: '/category/module-0-setup-orientation',
      },
      items: [
        'module-0-setup/environment',
        'module-0-setup/platform-orientation',
      ],
    },
    {
      type: 'category',
      label: 'Module 1 — Go Fundamentals',
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Module 1 — Go Fundamentals',
        description:
          'Weeks 2–4. The Go language from first principles: syntax, types, collections, structs, interfaces, errors, packages, and generics — taught with the idioms the platform expects from day one.',
        slug: '/category/module-1-go-fundamentals',
      },
      items: [
        'module-1-go-fundamentals/syntax-variables-types',
        'module-1-go-fundamentals/control-flow-functions',
        'module-1-go-fundamentals/collections',
        'module-1-go-fundamentals/structs-methods',
        'module-1-go-fundamentals/pointers-memory-basics',
        'module-1-go-fundamentals/interfaces',
        'module-1-go-fundamentals/error-handling',
        'module-1-go-fundamentals/packages-modules',
        'module-1-go-fundamentals/generics',
        'module-1-go-fundamentals/reflection-and-when-not',
      ],
    },
    {
      type: 'category',
      label: 'Module 2 — Intermediate Go',
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Module 2 — Intermediate Go',
        description:
          'Weeks 5–7. Concurrency, context, dependency injection, configuration, logging, testing, and performance — the engineering layer between "knows Go" and "ships Go services".',
        slug: '/category/module-2-intermediate-go',
      },
      items: [
        'module-2-intermediate/concurrency',
        'module-2-intermediate/context',
        'module-2-intermediate/dependency-injection',
        'module-2-intermediate/configuration',
        'module-2-intermediate/logging',
        'module-2-intermediate/testing',
        'module-2-intermediate/benchmarking-profiling',
        'module-2-intermediate/memory-performance',
      ],
    },
    {
      type: 'category',
      label: 'Module 3 — Advanced Go & Service Engineering',
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Module 3 — Advanced Go & Service Engineering',
        description:
          'Weeks 8–10. Production microservices: project structure, HTTP and middleware, REST APIs, auth, databases, transactions, RabbitMQ, workers, distributed systems, observability, containers, and CI/CD. Exercises run against the real local stack from here on.',
        slug: '/category/module-3-advanced-go',
      },
      items: [
        'module-3-advanced/project-structure',
        'module-3-advanced/http-servers-middleware',
        'module-3-advanced/rest-api-development',
        'module-3-advanced/authn-authz',
        'module-3-advanced/database-patterns',
        'module-3-advanced/transactions',
        'module-3-advanced/event-driven-rabbitmq',
        'module-3-advanced/workers-cron',
        'module-3-advanced/distributed-systems',
        'module-3-advanced/observability',
        'module-3-advanced/containers-kubernetes',
        'module-3-advanced/cicd',
      ],
    },
    {
      type: 'category',
      label: 'Module 4 — The DX Platform',
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Module 4 — The DX Platform',
        description:
          'Weeks 11–12. Everything specific to the Data Exchange: architecture, repository layout, the dx-common-go shared library, service anatomy, coding standards, security model, testing strategy, and deployment.',
        slug: '/category/module-4-the-dx-platform',
      },
      items: [
        'module-4-platform/architecture-deep-dive',
        'module-4-platform/repo-structure-workflow',
        'module-4-platform/dx-common-go-tour',
        'module-4-platform/service-anatomy',
        'module-4-platform/standards-checklist',
        'module-4-platform/security-model',
        'module-4-platform/testing-strategy',
        'module-4-platform/deployment',
      ],
    },
    {
      type: 'category',
      label: 'Module 5 — Capstone & First Contribution',
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Module 5 — Capstone & First Contribution',
        description:
          'The finish line: build a complete DX-style service with dx-common-go, then land your first real pull request on the platform.',
        slug: '/category/module-5-capstone',
      },
      items: ['capstone/capstone-service', 'capstone/first-contribution'],
    },
  ],
};

export default sidebars;
