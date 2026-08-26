import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'CDPG Go Developer Portal',
  tagline: 'Architecture, service engineering, and Go learning for the CDPG platform',
  favicon: 'img/favicon.svg',

  future: {
    v4: true,
  },

  url: 'https://datakaveri.github.io',
  baseUrl: '/go-learning/',

  organizationName: 'datakaveri',
  projectName: 'go-learning',

  onBrokenLinks: 'throw',

  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: 'throw',
    },
  },

  themes: ['@docusaurus/theme-mermaid'],

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/datakaveri/go-learning/tree/main/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/iudx-novo-social-card.png',
    colorMode: {
      defaultMode: 'light',
      respectPrefersColorScheme: true,
    },
    // Mermaid restyled to the CDPG palette (same values as cdpg-docs) so the
    // curriculum diagrams share the design language of the main docs site.
    mermaid: {
      theme: {light: 'base', dark: 'base'},
      options: {
        fontFamily: "Poppins, 'Segoe UI', system-ui, sans-serif",
        themeVariables: {
          fontSize: '14px',
          primaryColor: '#eef1f8',
          primaryTextColor: '#1f3569',
          primaryBorderColor: '#8794b8',
          secondaryColor: '#fdf0e3',
          tertiaryColor: '#f5f6f9',
          lineColor: '#8794b8',
          signalColor: '#8794b8',
          signalTextColor: '#5a6680',
          actorBkg: '#eef1f8',
          actorBorder: '#1f3569',
          actorTextColor: '#1f3569',
          labelBoxBkgColor: '#fdf0e3',
          labelBoxBorderColor: '#f57e20',
          labelTextColor: '#1f3569',
          loopTextColor: '#6a7488',
          noteBkgColor: '#fff4e8',
          noteBorderColor: '#f57e20',
          noteTextColor: '#1f3569',
          activationBkgColor: '#fdf0e3',
          activationBorderColor: '#f57e20',
          clusterBkg: '#f5f6f9',
          clusterBorder: '#d9deeb',
          edgeLabelBackground: '#ffffff',
        },
      },
    },
    navbar: {
      title: 'CDPG Go',
      logo: {
        alt: 'DX — Data Exchange Go Learning Path',
        src: 'img/logo.svg',
        srcDark: 'img/logo-dark.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'developerGuideSidebar',
          position: 'left',
          label: 'Developer Guide',
        },
        {
          type: 'docSidebar',
          sidebarId: 'curriculumSidebar',
          position: 'left',
          label: 'Go Curriculum',
        },
        {to: '/new-service/quick-start', label: 'New Service', position: 'left'},
        {
          href: 'https://datakaveri.github.io/cdpg-docs/',
          label: 'Platform Docs',
          position: 'right',
        },
        {
          href: 'https://datakaveri.github.io/dx-common-go-docs/',
          label: 'Go SDK',
          position: 'right',
        },
        {
          href: 'https://github.com/datakaveri/go-learning',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Learn',
          items: [
            {label: 'Start Here', to: '/start/learning-paths'},
            {label: 'New Service Tutorial', to: '/new-service/quick-start'},
            {label: 'Service Review Scorecard', to: '/standards/review-scorecard'},
          ],
        },
        {
          title: 'Go Resources',
          items: [
            {label: 'A Tour of Go', href: 'https://go.dev/tour/'},
            {label: 'Effective Go', href: 'https://go.dev/doc/effective_go'},
            {label: 'Go by Example', href: 'https://gobyexample.com/'},
          ],
        },
        {
          title: 'Platform',
          items: [
            {label: 'Architecture orientation', to: '/architecture/platform-orientation'},
            {label: 'Gateway integration', to: '/integrations/gateway'},
            {label: 'Platform docs', href: 'https://datakaveri.github.io/cdpg-docs/'},
            {label: 'dx-common-go docs', href: 'https://datakaveri.github.io/dx-common-go-docs/'},
            {label: 'GitHub', href: 'https://github.com/datakaveri'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Centre for Data for Public Good (CDPG), IISc. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'json', 'go', 'yaml', 'docker', 'sql'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
