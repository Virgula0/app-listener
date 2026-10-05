# Graph Report - app-listener  (2026-10-05)

## Corpus Check
- 314 files · ~1,489,203 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4166 nodes · 16535 edges · 171 communities (139 shown, 32 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1218 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c4c63835`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- runUpdate
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- Hash
- .objs
- User
- runUninstall
- highlight_test.go
- reloadOnce
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- Config
- usecase/daemon_test.go
- Resource
- inodeChain
- Candidate
- guard.bpf.c
- Run
- NetEventType
- SanitizeText
- buildOneGuard
- SelectOrphans
- serve.go
- engine
- StatFile
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- runNetworkMonitor
- GuardEvent
- guarded_ancestor_within_limit
- codeedit.go
- TempBinary
- absPath
- .write
- main_test.go
- os.FileMode
- Monitor
- planTaintOwners
- IntegrationSuite
- filevault.go
- RefreshSection
- IntegrationSuite
- Vault
- github.com/charmbracelet/bubbletea.Model
- IntegrationSuite
- os.File
- controlServer
- runForward
- EventType
- TrustGuard
- liftSuperseded
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- watchSet
- time.Duration
- github.com/spf13/cobra.Command
- runDaemon
- copyRegularAt
- bufio.Reader
- netModel
- btrfsLayoutFrom
- filevault_test.go
- KernelDev
- planUpdaters
- process_vm_readv.c
- TrustGuard
- systemd.go
- IntegrationSuite
- .step
- sync.Mutex
- net_tester/main.go
- render.go
- GuardInodeKey
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- changelog.go
- guardModel
- preload_wait.c
- .syncSetsLocked
- classifySupportError
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- TestMonitorUseCaseLifecycle
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- MonitorUseCase
- Build & sign release assets (reusable workflow)
- runOfflineEdit
- runMonitor
- Vault
- inode_dev
- app-listener Terminal Demo Recording
- CandidateDir
- .setInspectors
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- syncMap
- bpfLsmListed
- WebSocket /ws client
- find.go
- writeWithin
- install.sh
- .resyncLink
- launch_rule_at
- make test-integration (rootful Docker suite)
- revertSSHAgents
- supersede_probe.c
- github.com/cilium/ebpf.Map
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- IntegrationSuite
- runNetworkGuard
- vaultFS
- Guard
- tempRow
- editorModel
- configedit_test.go
- OpenSystemPlaced
- ensureMasterKey
- readKeyFrom
- parseGuardEvent
- SuperblockDev
- find_test.go
- DiscoverInfraBinaries
- launchKey
- guardLSMHooks
- Vault
- graphify-refresh.sh
- IntegrationSuite
- diffModel
- ParseEventsFlag
- .runEditorHarness
- ResolveTempBinary
- copyXattrs
- classCacheKey

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `IntegrationSuite` - 75 edges
5. `absPath()` - 74 edges
6. `IntegrationSuite` - 66 edges
7. `writeConfig()` - 65 edges
8. `fileEditModel` - 60 edges
9. `User` - 56 edges
10. `guardUnitTest` - 52 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (171 total, 32 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (120): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+112 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "wipe"
Cohesion: 0.11
Nodes (20): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr() (+12 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (69): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+61 more)

### Community 5 - "runUpdate"
Cohesion: 0.07
Nodes (37): changelogText(), TestChangelogText(), applyUpdate(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), downloadFile(), downloadReleaseFiles() (+29 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "rcUser"
Cohesion: 0.14
Nodes (32): revertBunLaunchers(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv() (+24 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (19): addErrKind, guardUnitTest, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+11 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "Hash"
Cohesion: 0.13
Nodes (29): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+21 more)

### Community 12 - "User"
Cohesion: 0.10
Nodes (32): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), pathCovered(), restoreDaemonAfterDiffAbort() (+24 more)

### Community 13 - "runUninstall"
Cohesion: 0.09
Nodes (31): runDiffCatalog(), parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), restoreBackups(), offerBackupCleanup() (+23 more)

### Community 14 - "highlight_test.go"
Cohesion: 0.15
Nodes (16): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 15 - "reloadOnce"
Cohesion: 0.11
Nodes (24): backingDeviceUnion(), buildConcurrency(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), startCatalogRefresh(), startGuardedDaemon() (+16 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (55): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf() (+47 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "Config"
Cohesion: 0.07
Nodes (46): ServeConfig, sameConfig(), newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), inspectorPaths(), resolveInspectors() (+38 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (53): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+45 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): fakeDaemon, Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "inodeChain"
Cohesion: 0.18
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 24 - "Candidate"
Cohesion: 0.18
Nodes (17): selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories(), groupCandidates(), pickDirectories(), pickFromCandidates(), pickUsers(), TestGroupCandidates() (+9 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (44): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector(), discover_guarded_parent() (+36 more)

### Community 26 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 27 - "NetEventType"
Cohesion: 0.07
Nodes (21): NetBpfEvent, NetEvent, NetEventType, foreignPinnedLink(), Cstr(), FormatAddr(), NetEventTypes(), Ntohs() (+13 more)

### Community 28 - "SanitizeText"
Cohesion: 0.06
Nodes (28): newBinaryVetter(), openBinaryVetter(), ownerPaths(), grantNewBinaries(), runHeadless(), runHeadless(), candidates(), confirmable() (+20 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.08
Nodes (37): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardTUI() (+29 more)

### Community 30 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 31 - "serve.go"
Cohesion: 0.06
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 33 - "StatFile"
Cohesion: 0.18
Nodes (7): refusedReplacements, ReplacementCheck, TrustGuard, SetReplacementCheck(), TestReSyncBinaries_RefusedReplacementKeepsOldKey(), TestReSyncBinaries_ReplacementJudgedAgainstPrunedKey(), StatFile()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.08
Nodes (45): AtomicWriteAt(), OpenRegularNoFollow(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf() (+37 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.09
Nodes (33): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+25 more)

### Community 36 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (9): ParseNetEventsFlag(), runNetworkMonitor(), runTUI(), NetworkMonitorRepository, NewNetModel(), assertViewFits(), TestViewsFitReportedTerminalSize(), NetworkMonitorUseCase (+1 more)

### Community 38 - "GuardEvent"
Cohesion: 0.12
Nodes (7): holdHeadless(), HeadlessLine(), runGuardHeadless(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial(), GuardUseCase

### Community 39 - "guarded_ancestor_within_limit"
Cohesion: 0.24
Nodes (19): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename() (+11 more)

### Community 40 - "codeedit.go"
Cohesion: 0.19
Nodes (25): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+17 more)

### Community 41 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+8 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - ".write"
Cohesion: 0.16
Nodes (13): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+5 more)

### Community 44 - "main_test.go"
Cohesion: 0.07
Nodes (30): CheckBPFLSM(), runBPFCheck(), launchCase, guardBinaryFlag(), infraContainerPaths(), TestIntegrationSuite(), TestMain(), IntegrationSuite (+22 more)

### Community 45 - "os.FileMode"
Cohesion: 0.10
Nodes (28): TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent() (+20 more)

### Community 46 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 48 - "planTaintOwners"
Cohesion: 0.17
Nodes (13): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), isSubset(), planMemberRows(), planTaintOwners(), setKeyOf() (+5 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "filevault.go"
Cohesion: 0.23
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 51 - "RefreshSection"
Cohesion: 0.11
Nodes (20): logRefreshChange(), TestLogRefreshChange(), updateCatalogConfig(), RefreshOptions, SectionScan, TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff() (+12 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): IntegrationSuite, exploitTest, IntegrationSuite, newEventTypes(), tailLast(), IntegrationSuite, editorHarness

### Community 53 - "Vault"
Cohesion: 0.15
Nodes (15): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+7 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.13
Nodes (6): renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 56 - "os.File"
Cohesion: 0.17
Nodes (16): checkFreshDir(), fileID(), trustManager, openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir() (+8 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "runForward"
Cohesion: 0.18
Nodes (14): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), chooseResources() (+6 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (14): systemRule(), BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+6 more)

### Community 60 - "TrustGuard"
Cohesion: 0.14
Nodes (5): trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "liftSuperseded"
Cohesion: 0.12
Nodes (13): binaryVerifyState, sharedHashEntry, hashBinaryShared(), dropSupersededLocked(), exeGone(), Guard, SupersededBinary, liftSuperseded() (+5 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+9 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+20 more)

### Community 66 - "watchSet"
Cohesion: 0.06
Nodes (36): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+28 more)

### Community 67 - "time.Duration"
Cohesion: 0.13
Nodes (8): drain(), newCatalogRefresher(), readEvents(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), catalogRefresher, debouncer

### Community 68 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (33): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+25 more)

### Community 69 - "runDaemon"
Cohesion: 0.15
Nodes (14): awaitStartupOrSignal(), catchLifecycleSignals(), confirmMasterKeyOverwrite(), notifySystemdReady(), runDaemon(), runGenKey(), runOneShotMode(), startGuardedDaemonAbortable() (+6 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.14
Nodes (13): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+5 more)

### Community 73 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+18 more)

### Community 76 - "KernelDev"
Cohesion: 0.18
Nodes (13): BinaryStat, binaryStatOf(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs() (+5 more)

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 81 - "systemd.go"
Cohesion: 0.11
Nodes (34): reloadDaemonForPasswordChange(), buildBinaryIfNeeded(), installBinaryOnly(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestSoftenAutomatedRefreshErr(), mustCwd(), revertSystemFiles() (+26 more)

### Community 83 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 84 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "GuardInodeKey"
Cohesion: 0.27
Nodes (4): TestNearestRootClaims(), TestTieRank(), nearestRootClaims(), tieRank()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "changelog.go"
Cohesion: 0.22
Nodes (7): newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 92 - "guardModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 94 - ".syncSetsLocked"
Cohesion: 0.27
Nodes (4): engine, intersectAllows(), setPaths(), TestIntersectAllows()

### Community 95 - "classifySupportError"
Cohesion: 0.18
Nodes (10): classifySetupError(), classifySupportError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), TestClassifySupportErrorEncryptionNotEnabled(), TestClassifySupportErrorEncryptionNotEnabledF2fs(), TestClassifySupportErrorGeneric() (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "liveSession"
Cohesion: 0.22
Nodes (7): parseErr(), displayEntries(), holdForward(), pingInterval(), liveSession, RunFileEditorSession(), RunSession()

### Community 100 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.33
Nodes (8): NewGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle(), fakeNetworkGuardRepo

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "MonitorUseCase"
Cohesion: 0.20
Nodes (5): runHeadless(), MonitorRepository, MonitorUseCase, NewMonitorUseCase(), TestMonitorUseCaseStartError()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runOfflineEdit"
Cohesion: 0.07
Nodes (25): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+17 more)

### Community 110 - "runMonitor"
Cohesion: 0.10
Nodes (24): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair() (+16 more)

### Community 111 - "Vault"
Cohesion: 0.18
Nodes (10): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), copyTreeRoot(), CopyTreeWithProgress() (+2 more)

### Community 112 - "inode_dev"
Cohesion: 0.31
Nodes (6): ctx_ptr(), inode_dev(), sb_dev(), guard_bprm_committed(), guard_inode_free(), trust_inode_free()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "CandidateDir"
Cohesion: 0.05
Nodes (35): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+27 more)

### Community 115 - ".setInspectors"
Cohesion: 0.33
Nodes (4): engine, Inspector, TrustGuard, SetInspectors()

### Community 116 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.26
Nodes (3): column(), Finder, overlay()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "syncMap"
Cohesion: 0.29
Nodes (6): syncMap(), OpenConfined(), ResolveConfined(), StatConfined(), TestResolveConfined(), TestResolveConfinedRootOwnedParents()

### Community 121 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "find.go"
Cohesion: 0.27
Nodes (7): findButton, span, findAll(), lowerRunes(), newFindInput(), runesEqual(), TestFindAllSmartCase()

### Community 124 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - ".resyncLink"
Cohesion: 0.39
Nodes (3): engine, Guard, supersedeUnnamed()

### Community 127 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "revertSSHAgents"
Cohesion: 0.25
Nodes (9): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers(), TestIsInstallerSSHAgentUnit(), TestRemoveSSHAgentUnit() (+1 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "github.com/cilium/ebpf.Map"
Cohesion: 0.25
Nodes (5): deleteInoKeys(), deleteResKeys(), syncSetMap(), syncU32Map(), deleteKeys()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 144 - "runNetworkGuard"
Cohesion: 0.10
Nodes (17): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetGuardEvent (+9 more)

### Community 145 - "vaultFS"
Cohesion: 0.36
Nodes (3): applyNewFileMeta(), listEntries(), vaultFS

### Community 146 - "Guard"
Cohesion: 0.10
Nodes (17): deferredBinary, VettedInode, BinaryEntry, SharedPinDegraded(), canonicalBinaryPath(), ConfinedEntry(), confinedEntry(), eventMask() (+9 more)

### Community 147 - "tempRow"
Cohesion: 0.20
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust()

### Community 148 - "editorModel"
Cohesion: 0.33
Nodes (4): nextConfig(), editorModel, EditText(), newEditorModel()

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "OpenSystemPlaced"
Cohesion: 0.12
Nodes (16): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), fileKey() (+8 more)

### Community 151 - "ensureMasterKey"
Cohesion: 0.33
Nodes (6): ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), GenerateMasterKey(), MasterKeyExists()

### Community 152 - "readKeyFrom"
Cohesion: 0.33
Nodes (6): checkKeyLen(), readKeyFrom(), readOrCreateKey(), TestReadKeyFromBadLength(), TestReadKeyFromMissingFile(), TestReadKeyFromRoundTrip()

### Community 153 - "parseGuardEvent"
Cohesion: 0.22
Nodes (7): bpfGuardEvent, fsGateLabel(), parseGuardEvent(), processGateLabel(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 154 - "SuperblockDev"
Cohesion: 0.15
Nodes (12): checkBtrfsKeys(), entryOf(), BtrfsMounted(), FirstOnBtrfs(), mountinfoDev(), mountinfoHasFstype(), OnBtrfs(), SuperblockDev() (+4 more)

### Community 155 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 156 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 157 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 158 - "guardLSMHooks"
Cohesion: 0.40
Nodes (3): trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 163 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 165 - "ResolveTempBinary"
Cohesion: 0.20
Nodes (7): hashNow(), loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows(), TestResolveTempBinaryClassesRuntime(), ComputeBinaryEntryFile()

### Community 166 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

## Knowledge Gaps
- **52 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 375 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `StatFile`, `DaemonEvent`, `GuardEvent`, `.objs`, `reloadOnce`, `liftSuperseded`, `sync.Mutex`, `Resource`, `os.File`, `buildOneGuard`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `monitorUnitTest` connect `Monitor` to `daemon/daemon.go`, `IntegrationSuite`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09350991858731797 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024523523923763828 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._