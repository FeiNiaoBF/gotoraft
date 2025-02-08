"use client"

import { PageLayout } from "@/components/page-layout"
import { BackToHome } from "@/components/back-to-home"
import { useLanguage } from "@/contexts/language-context"

export default function ConfigurePage() {
  const { t } = useLanguage()

  return (
    <PageLayout title={t("configureTitle")} description={t("Adjust the parameters of your Raft cluster simulation.")}>
      <div className="mb-6">
        <BackToHome />
      </div>
      <div className="bg-slate-800/50 p-6 rounded-lg">
        <p className="text-slate-300">{t("Node configuration controls will be implemented here.")}</p>
      </div>
    </PageLayout>
  )
}

