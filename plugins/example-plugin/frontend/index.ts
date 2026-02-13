/**
 * # Example Plugin - Frontend
 *
 * Panel.go için örnek frontend plugin implementasyonu.
 * Custom field component'i gösterir.
 *
 * ## Kullanım
 * ```typescript
 * // main.tsx veya plugin loader'da
 * import { ExamplePlugin } from './plugins/example-plugin/frontend';
 * import { pluginRegistry } from '@/plugins/PluginRegistry';
 *
 * pluginRegistry.register(ExamplePlugin);
 * ```
 */

import { Plugin } from '@/plugins/types';
import { ExampleField } from './fields/ExampleField';

/**
 * Example Plugin
 *
 * Custom field ve widget örneği içerir.
 */
export const ExamplePlugin: Plugin = {
  name: 'example-plugin',
  version: '1.0.0',
  description: 'Example plugin demonstrating plugin system',
  author: 'Panel.go Team',

  // Custom field'lar
  fields: [
    {
      type: 'example-field',
      component: ExampleField,
    },
  ],

  // Plugin initialization
  init: async () => {
    console.log('ExamplePlugin initialized');
  },

  // Plugin cleanup
  cleanup: async () => {
    console.log('ExamplePlugin cleaned up');
  },
};

// Auto-register (import in plugins/index.ts)
