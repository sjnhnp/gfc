import { ReadFile, WriteFile, SaveFileDialog, OpenFileDialog, ReadExternalFile, WriteExternalFile } from '@/bridge'
import { useProfilesStore, useSubscribesStore, useRulesetsStore, usePluginsStore } from '@/stores'
import { ProfilesFilePath, SubscribesFilePath, RulesetsFilePath, PluginsFilePath } from '@/constant'

export interface ExportData {
    version: string
    exportedAt: string
    profiles: any[]
    subscribes: any[]
    rulesets: any[]
    plugins: any[]
}

/**
 * Export all configurations to a JSON file
 */
export const exportAllConfigs = async (): Promise<boolean> => {
    const profilesStore = useProfilesStore()
    const subscribesStore = useSubscribesStore()
    const rulesetsStore = useRulesetsStore()
    const pluginsStore = usePluginsStore()

    const data: ExportData = {
        version: '1.0.0',
        exportedAt: new Date().toISOString(),
        profiles: profilesStore.profiles,
        subscribes: subscribesStore.subscribes,
        rulesets: rulesetsStore.rulesets,
        plugins: pluginsStore.plugins,
    }

    const filePath = await SaveFileDialog(
        'Export Configuration',
        `gfc-config-${new Date().toISOString().split('T')[0]}.json`,
        'JSON Files:*.json|All Files:*.*'
    )

    if (!filePath) {
        return false // User cancelled
    }

    await WriteExternalFile(filePath, JSON.stringify(data, null, 2))
    return true
}

/**
 * Import configurations from a JSON file
 */
export const importAllConfigs = async (options: {
    profiles?: boolean
    subscribes?: boolean
    rulesets?: boolean
    plugins?: boolean
    merge?: boolean // If true, merge with existing; if false, replace
} = {}): Promise<boolean> => {
    const {
        profiles = true,
        subscribes = true,
        rulesets = true,
        plugins = true,
        merge = false
    } = options

    const filePath = await OpenFileDialog(
        'Import Configuration',
        'JSON Files:*.json|All Files:*.*'
    )

    if (!filePath) {
        return false // User cancelled
    }

    const content = await ReadExternalFile(filePath)
    const data: ExportData = JSON.parse(content)

    // Validate data structure
    if (!data.version || !data.exportedAt) {
        throw new Error('Invalid configuration file format')
    }

    const profilesStore = useProfilesStore()
    const subscribesStore = useSubscribesStore()
    const rulesetsStore = useRulesetsStore()
    const pluginsStore = usePluginsStore()

    if (profiles && data.profiles) {
        if (merge) {
            // Merge: add only if ID doesn't exist
            for (const profile of data.profiles) {
                if (!profilesStore.getProfileById(profile.id)) {
                    profilesStore.profiles.push(profile)
                }
            }
        } else {
            profilesStore.profiles.splice(0, profilesStore.profiles.length, ...data.profiles)
        }
        await profilesStore.saveProfiles()
    }

    if (subscribes && data.subscribes) {
        if (merge) {
            for (const sub of data.subscribes) {
                if (!subscribesStore.getSubscribeById(sub.id)) {
                    subscribesStore.subscribes.push(sub)
                }
            }
        } else {
            subscribesStore.subscribes.splice(0, subscribesStore.subscribes.length, ...data.subscribes)
        }
        await subscribesStore.saveSubscribes()
    }

    if (rulesets && data.rulesets) {
        if (merge) {
            for (const ruleset of data.rulesets) {
                if (!rulesetsStore.getRulesetById(ruleset.id)) {
                    rulesetsStore.rulesets.push(ruleset)
                }
            }
        } else {
            rulesetsStore.rulesets.splice(0, rulesetsStore.rulesets.length, ...data.rulesets)
        }
        await rulesetsStore.saveRulesets()
    }

    if (plugins && data.plugins) {
        if (merge) {
            for (const plugin of data.plugins) {
                if (!pluginsStore.getPluginById(plugin.id)) {
                    pluginsStore.plugins.push(plugin)
                }
            }
        } else {
            pluginsStore.plugins.splice(0, pluginsStore.plugins.length, ...data.plugins)
        }
        await pluginsStore.savePlugins()
    }

    return true
}
