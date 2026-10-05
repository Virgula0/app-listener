# Graph Report - app-listener  (2026-10-05)

## Corpus Check
- 312 files · ~1,484,987 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4117 nodes · 16324 edges · 157 communities (131 shown, 26 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1189 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7171dc5e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- update.go
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- Hash
- Guard
- Monitor
- uninstall.go
- install.go
- runDaemon
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- Config
- usecase/daemon_test.go
- Resource
- backups.go
- go_pkg_path_filepath
- guard.bpf.c
- Run
- NetEventType
- SanitizeText
- buildOneGuard
- buildTrustedSet
- serve.go
- engine
- os.File
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- runNetworkMonitor
- GuardEvent
- KernelDev
- codeedit.go
- highlight_test.go
- IntegrationSuite
- fscrypt/orphans.go
- IntegrationSuite
- safeio.go
- IntegrationSuite
- .syncSetsLocked
- IntegrationSuite
- runMonitor
- BinaryEntry
- wipe
- monitorUnitTest
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- catalogRefresher
- .serveRequest
- TestMonitorUseCaseLifecycle
- FileEvent
- TrustGuard
- supersedeMaps
- exe_supersede.h
- fakeGuardRepo
- string
- pinstate.go
- watchSet
- guarded_ancestor_within_limit
- NewDumpHook
- .add
- copytree.go
- github.com/spf13/cobra.Command
- bpfstats/main.go
- btrfsLayoutFrom
- tempRow
- writeWithin
- planUpdaters
- process_vm_readv.c
- github.com/cilium/ebpf.Map
- LibraryClosure
- compileGlobNames
- walkPlaced
- controlServer
- net_tester/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DiscoverForUsers
- time.Duration
- preload_wait.c
- netModel
- TempBinary
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- IntegrationSuite
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- vaultFS
- Build & sign release assets (reusable workflow)
- New
- parseConfig
- Vault
- inode_dev
- app-listener Terminal Demo Recording
- CandidateDir
- .setInspectors
- runForward
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- tempGrant
- bpfLsmListed
- WebSocket /ws client
- lockRootRecovering
- tempGrantSetup
- install.sh
- IntegrationSuite
- launch_rule_at
- make test-integration (rootful Docker suite)
- loadDaemonConfig
- supersede_probe.c
- fakeDaemon
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- .resyncLink
- ParseEventsFlag
- runNetworkGuard
- sealFileVaultWithMasterKey
- EventType
- TemporaryGrant
- unlockUnderGuard
- configedit_test.go
- launchKey
- guardLSMHooks
- lockAndDeprovision
- .checkCommSpoof
- .write
- Vault
- graphify-refresh.sh

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `IntegrationSuite` - 75 edges
5. `absPath()` - 73 edges
6. `IntegrationSuite` - 66 edges
7. `writeConfig()` - 65 edges
8. `fileEditModel` - 59 edges
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

## Communities (157 total, 26 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (31): Execute(), tempJournalRow, classCacheKey, launchRule, trustEvent, btrfsIoctlFsInfoArgs, attrOpCase, guardExploitTest (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (124): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+116 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.20
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+2 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): TestSetSectionWhitelistPreservesGroupStructure(), updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+71 more)

### Community 5 - "update.go"
Cohesion: 0.05
Nodes (52): decryptDirectories(), TestDecryptDirectoriesEmpty(), applyUpdate(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), downloadFile(), downloadReleaseFiles() (+44 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "shellrc.go"
Cohesion: 0.07
Nodes (69): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), catalogPatterns(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor() (+61 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "Hash"
Cohesion: 0.12
Nodes (30): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+22 more)

### Community 11 - "Guard"
Cohesion: 0.08
Nodes (7): binaryVerifyState, rootHandle, sharedHashEntry, openRoot(), SharedPinDegraded(), Guard, hashBinaryShared()

### Community 12 - "Monitor"
Cohesion: 0.12
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 13 - "uninstall.go"
Cohesion: 0.06
Nodes (43): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askDecrypt(), cleanOrphanedMetadata() (+35 more)

### Community 14 - "install.go"
Cohesion: 0.06
Nodes (61): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+53 more)

### Community 15 - "runDaemon"
Cohesion: 0.06
Nodes (42): CheckBPFLSM(), CheckEBPF(), awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), configureDaemonLogging() (+34 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.12
Nodes (7): eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.06
Nodes (52): logRefreshChange(), TestLogRefreshChange(), askEncryption(), LibraryBlock, RefreshOptions, Section, findSectionEnd(), findSectionStart() (+44 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "Config"
Cohesion: 0.09
Nodes (37): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), auditAfterEdit(), auditAfterEditWithConfig(), auditEntry() (+29 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.10
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "backups.go"
Cohesion: 0.11
Nodes (27): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), Delete() (+19 more)

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
Cohesion: 0.07
Nodes (24): newBinaryVetter(), openBinaryVetter(), ownerPaths(), candidates(), confirmable(), confirmAll(), describe(), run() (+16 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.10
Nodes (31): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardTUI() (+23 more)

### Community 30 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 31 - "serve.go"
Cohesion: 0.06
Nodes (36): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+28 more)

### Community 32 - "engine"
Cohesion: 0.11
Nodes (6): TestNearestRootClaims(), TestTieRank(), engine, nearestRootClaims(), pathWithin(), tieRank()

### Community 33 - "os.File"
Cohesion: 0.08
Nodes (28): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+20 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.09
Nodes (45): AtomicWriteAt(), OpenRegularNoFollow(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf() (+37 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.08
Nodes (38): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+30 more)

### Community 36 - "watchGroup"
Cohesion: 0.10
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "runNetworkMonitor"
Cohesion: 0.14
Nodes (9): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase (+1 more)

### Community 38 - "GuardEvent"
Cohesion: 0.07
Nodes (21): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, fsGateLabel(), GuardEvent, logBacklogDenial(), parseGuardEvent() (+13 more)

### Community 39 - "KernelDev"
Cohesion: 0.09
Nodes (29): fdStat, BinaryStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD() (+21 more)

### Community 40 - "codeedit.go"
Cohesion: 0.16
Nodes (28): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+20 more)

### Community 41 - "highlight_test.go"
Cohesion: 0.09
Nodes (25): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+17 more)

### Community 42 - "IntegrationSuite"
Cohesion: 0.10
Nodes (4): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite

### Community 43 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): IntegrationSuite, exploitTest, IntegrationSuite, newEventTypes(), tailLast(), IntegrationSuite

### Community 45 - "safeio.go"
Cohesion: 0.07
Nodes (37): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+29 more)

### Community 46 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 48 - ".syncSetsLocked"
Cohesion: 0.11
Nodes (18): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, intersectAllows(), isSubset(), planMemberRows() (+10 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "runMonitor"
Cohesion: 0.10
Nodes (21): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+13 more)

### Community 51 - "BinaryEntry"
Cohesion: 0.12
Nodes (17): deferredBinary, VettedInode, BinaryEntry, BinariesSummary(), canonicalBinaryPath(), ConfinedEntry(), confinedEntry(), patternMatches() (+9 more)

### Community 52 - "wipe"
Cohesion: 0.23
Nodes (17): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+9 more)

### Community 53 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (9): diffModel, editorModel, clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, serveTestModel, sessionModel (+1 more)

### Community 55 - "absPath"
Cohesion: 0.10
Nodes (7): IntegrationSuite, IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "catalogRefresher"
Cohesion: 0.19
Nodes (6): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher

### Community 57 - ".serveRequest"
Cohesion: 0.14
Nodes (15): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), streamEvents() (+7 more)

### Community 58 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.33
Nodes (8): NewGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle(), fakeNetworkGuardRepo

### Community 59 - "FileEvent"
Cohesion: 0.10
Nodes (8): BpfEvent, eventUnitTest, FileEvent, Run(), DecodeBpfEvent(), newFakeMonitorRepo(), editorHarness, fakeMonitorRepo

### Community 60 - "TrustGuard"
Cohesion: 0.15
Nodes (6): trustHook, trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 61 - "supersedeMaps"
Cohesion: 0.13
Nodes (9): dropSupersededLocked(), exeGone(), Guard, SupersededBinary, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop() (+1 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+9 more)

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "watchSet"
Cohesion: 0.16
Nodes (11): addSystemBinary(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify, watchDir (+3 more)

### Community 67 - "guarded_ancestor_within_limit"
Cohesion: 0.24
Nodes (19): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename() (+11 more)

### Community 68 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 69 - ".add"
Cohesion: 0.19
Nodes (4): TemporaryAccess, tempJournalRows(), fakeTempGrant, tempTrace

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 73 - "bpfstats/main.go"
Cohesion: 0.16
Nodes (18): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), NewTrustGuard(), TestTrustGuardStartFailsWhenLibraryAllowlistHookCannotAttach(), inertProgram(), kernelHasFunc() (+10 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "tempRow"
Cohesion: 0.27
Nodes (5): exeRow, tempMaskOp, tempRow, Guard, planTempBlock()

### Community 76 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 78 - "planUpdaters"
Cohesion: 0.10
Nodes (23): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+15 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 80 - "github.com/cilium/ebpf.Map"
Cohesion: 0.18
Nodes (6): deleteInoKeys(), deleteResKeys(), syncU32Map(), syncMap(), deleteKeys(), TrustGuard

### Community 81 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+8 more)

### Community 82 - "compileGlobNames"
Cohesion: 0.16
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), TrustGuard, resolveBits(), resolveChildren(), statKey() (+2 more)

### Community 83 - "walkPlaced"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 84 - "controlServer"
Cohesion: 0.13
Nodes (8): controlManager, controlServer, controlManager, controlServer, newControlManager(), startControlManager(), startControlServer(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.18
Nodes (12): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), repeatSpaces(), setGutter() (+4 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DiscoverForUsers"
Cohesion: 0.22
Nodes (11): configPaths(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), cleanOrphanedFscrypt(), isInsidePath(), Discover() (+3 more)

### Community 92 - "time.Duration"
Cohesion: 0.21
Nodes (4): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 94 - "netModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 95 - "TempBinary"
Cohesion: 0.26
Nodes (11): TempBinary, TempRule, ResolveTempBinary(), closeTempBinaries(), daemonUseCase, logTempGrant(), newExeTap(), planTempGrants() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

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

### Community 107 - "vaultFS"
Cohesion: 0.36
Nodes (3): applyNewFileMeta(), listEntries(), vaultFS

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "New"
Cohesion: 0.22
Nodes (15): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+7 more)

### Community 110 - "parseConfig"
Cohesion: 0.17
Nodes (11): parseConfig(), checkBtrfsKeys(), ensureBtrfsLayout(), entryOf(), FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev() (+3 more)

### Community 111 - "Vault"
Cohesion: 0.26
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata()

### Community 112 - "inode_dev"
Cohesion: 0.31
Nodes (6): ctx_ptr(), inode_dev(), sb_dev(), guard_bprm_committed(), guard_inode_free(), trust_inode_free()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "CandidateDir"
Cohesion: 0.07
Nodes (25): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+17 more)

### Community 115 - ".setInspectors"
Cohesion: 0.46
Nodes (3): engine, Inspector, SetInspectors()

### Community 116 - "runForward"
Cohesion: 0.20
Nodes (11): chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents(), pingInterval() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "tempGrant"
Cohesion: 0.21
Nodes (4): tempGrant, loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 121 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "lockRootRecovering"
Cohesion: 0.25
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 124 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 127 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "loadDaemonConfig"
Cohesion: 0.25
Nodes (8): loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestResolveConfinedRootOwnedParents(), TestOpenSystemPlacedRefusesUserHome()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - ".resyncLink"
Cohesion: 0.39
Nodes (3): engine, Guard, supersedeUnnamed()

### Community 143 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 144 - "runNetworkGuard"
Cohesion: 0.10
Nodes (18): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode (+10 more)

### Community 145 - "sealFileVaultWithMasterKey"
Cohesion: 0.23
Nodes (12): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape() (+4 more)

### Community 146 - "EventType"
Cohesion: 0.33
Nodes (5): systemRule(), EventType, parseEvents(), TestAddBinaryEventsReadOnlyRejectsRestriction(), ParseEventType()

### Community 147 - "TemporaryGrant"
Cohesion: 0.40
Nodes (3): TemporaryGrant, applyTempGrants(), revokeTempGrants()

### Community 148 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 152 - "lockAndDeprovision"
Cohesion: 0.27
Nodes (7): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 156 - ".write"
Cohesion: 0.22
Nodes (8): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), writeFileInPlace()

## Knowledge Gaps
- **52 isolated node(s):** `github.com/Virgula0/app-listener`, `hint`, `tempJournalRow`, `launchRule`, `trustEvent` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 373 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `os.File`, `DaemonEvent`, `fakeDaemon`, `GuardEvent`, `runDaemon`, `BinaryEntry`, `unlockUnderGuard`, `controlServer`, `.checkCommSpoof`, `buildOneGuard`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `hint`, `tempJournalRow` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10014665828619317 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.024606299212598427 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08982456140350877 - nodes in this community are weakly interconnected._