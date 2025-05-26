class ConfigService {
  constructor() {
    this.config = null;
    this.loaded = false;
  }

  async loadConfig() {
    if (this.loaded) {
      return this.config;
    }

    this.config = {
      backendUrl: process.env.VUE_APP_BACKEND_URL || window.location.origin
    };
    this.loaded = true;
    return this.config;
  }

  getBackendUrl() {
    if (!this.config) {
      throw new Error('Config not loaded. Call loadConfig() first.');
    }
    return this.config.backendUrl;
  }
}

export default new ConfigService();