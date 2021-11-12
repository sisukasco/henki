module.exports = {
    preset: 'ts-jest',
    testEnvironment: "node",
    testMatch:["**/tests/**/*.[jt]s"],
    testTimeout: 30000,
    "globals": {
        "test_config": {
            endpoint: "http://local.dockform.com:3131",
            aud: "local.dockform.com",
            emails: {
                endpoint: "http://local.dockform.com:8025"
            }
        }
    }
}