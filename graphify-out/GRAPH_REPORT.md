# Graph Report - app-listener  (2026-10-16)

## Corpus Check
- 379 files · ~1,572,212 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4761 nodes · 19980 edges · 187 communities (161 shown, 26 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1669 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1a81736e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- daemon_multicall_test.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- BinaryEntry
- liveSession
- go_pkg_path_filepath
- Config
- install.go
- Candidate
- .startGuardStd
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- newFakeVault
- Resource
- auth.go
- IntegrationSuite
- guard.bpf.c
- NetEventType
- highlight_test.go
- LibraryClosure
- IntegrationSuite
- Vault
- serve.go
- bpfstats/main.go
- changelog_test.go
- fileedit_model_test.go
- Ledger
- uninstall.go
- Monitor
- sync.Mutex
- ConfinedEntry
- fileedit.go
- GenerateConf
- IntegrationSuite
- monitorUnitTest
- newWatchPattern
- update_test.go
- Read
- stdio
- watchSet
- .FilterExistingWhitelist
- planUpdaters
- codeedit.go
- SelectOrphans
- watchGroup
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- unit.go
- absPath
- GuardInodeKey
- github.com/charmbracelet/bubbles/textarea.Model
- NetworkMonitor
- fileedit_symlink_test.go
- exe_supersede.h
- GuardEvent
- expandPlaceholders
- pinstate.go
- gui/session.go
- inode_dev
- Guard
- check_watched
- fcntl
- copyRegularAt
- .judgeProcess
- globBuilder
- btrfsLayoutFrom
- RefreshSection
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- process_vm_readv.c
- net_tester/main.go
- hashCache
- CandidateDir
- wizard.go
- IntegrationSuite
- TestMonitorUseCaseLifecycle
- tempGrantSetup
- fakeGuardRepo
- testKeyPair
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- catalogRefresher
- time.Time
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- NewDumpHook
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- Well-known bypass classes
- time.Duration
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- runNetworkGuard
- Build & sign release assets (reusable workflow)
- MonitorUseCase
- github.com/stretchr/testify/suite.Suite
- KernelDev
- find.go
- app-listener Terminal Demo Recording
- IntegrationSuite
- EventType
- hashFileShared
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- compileGlobNames
- .populateWatchInodes
- WebSocket /ws client
- markReady
- launch_probe.c
- install.sh
- mc_tag.h
- inode_dev.h
- make test-integration (rootful Docker suite)
- mountVouchesOwnership
- supersede_probe.c
- guardLSMHooks
- trace-app-libs.sh
- IntegrationSuite
- writeWithin
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- .exeKeys
- vouchedDevs
- mountTable
- buildOneGuard
- ldConfParser
- render.go
- launch_rule_at
- newTestVetter
- filevault.go
- IntegrationSuite
- daemon: permanent protection with encryption at rest
- integrationtests/guard_test.go
- TestMain
- io.Reader
- find_test.go
- bpfLsmListed
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- configedit_test.go
- inodeChain
- Vault
- BinaryStat
- Limitations
- Troubleshooting
- mountsUnitTest
- Daemon configuration
- edit-protected
- Compatibility
- How it works
- openBunTmp
- pathCache
- SanitizeText
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .eventTypeAllowed
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 80 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `fileEditModel` - 60 edges
10. `engine` - 58 edges

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

## Communities (187 total, 26 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (24): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (109): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags() (+101 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.05
Nodes (90): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions() (+82 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.21
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.08
Nodes (18): deferredBinary, BinaryEntry, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), BinariesSummary(), canonicalBinaryPath(), confinedEntry() (+10 more)

### Community 11 - "liveSession"
Cohesion: 0.15
Nodes (15): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+7 more)

### Community 12 - "go_pkg_path_filepath"
Cohesion: 0.34
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 13 - "Config"
Cohesion: 0.06
Nodes (52): newBinaryVetter(), openBinaryVetter(), makeReloadHandler(), newDaemonModel(), noteReloadedSettings(), runTUI(), startCatalogRefresh(), inspectorPaths() (+44 more)

### Community 14 - "install.go"
Cohesion: 0.12
Nodes (25): confirmMasterKeyOverwrite(), runGenKey(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), ensureInstalledBinary(), ensureMasterKey() (+17 more)

### Community 15 - "Candidate"
Cohesion: 0.17
Nodes (20): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+12 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.12
Nodes (9): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent() (+1 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "newFakeVault"
Cohesion: 0.13
Nodes (47): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): fakeDaemon, Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.14
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class(), evict_inode_from_guard() (+45 more)

### Community 26 - "NetEventType"
Cohesion: 0.09
Nodes (12): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), Ntohs(), ParseNetEventType(), reasonText() (+4 more)

### Community 27 - "highlight_test.go"
Cohesion: 0.15
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 28 - "LibraryClosure"
Cohesion: 0.16
Nodes (13): defaultLibDirs(), elfInterp(), fileExists(), isScript(), LdSoPreloadPaths(), LibraryClosure(), resolveSoname(), TestResolveLibraryClosureRejectsUserWritableOriginDir() (+5 more)

### Community 29 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 30 - "Vault"
Cohesion: 0.07
Nodes (30): encryptDirectories(), secureResources(), verifyEncryptionState(), resolveFilesystemPrereqs(), deprovisionKind, harnessSuite, migrationTarget, isFileVaultCiphertext() (+22 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "changelog_test.go"
Cohesion: 0.17
Nodes (11): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+3 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.25
Nodes (17): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+9 more)

### Community 35 - "Ledger"
Cohesion: 0.10
Nodes (13): logRefreshChange(), TestLogRefreshChange(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run() (+5 more)

### Community 36 - "uninstall.go"
Cohesion: 0.05
Nodes (63): reloadDaemonForPasswordChange(), preflightDeployedBinary(), confirmRunPrereq(), installConfig(), revertSystemFiles(), applyUpdate(), checkKeyLen(), classifySetupError() (+55 more)

### Community 37 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 38 - "sync.Mutex"
Cohesion: 0.08
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+15 more)

### Community 39 - "ConfinedEntry"
Cohesion: 0.27
Nodes (20): ConfinedEntry(), HashStats(), newVerifyState(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash() (+12 more)

### Community 40 - "fileedit.go"
Cohesion: 0.09
Nodes (26): checkFreshDir(), openDirNoFollow(), renewDir(), AtomicWriteAt(), OpenRegularNoFollow(), applyNewFileMeta(), defaultNewFileMode(), fileModeToUnixMode() (+18 more)

### Community 41 - "GenerateConf"
Cohesion: 0.08
Nodes (38): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+30 more)

### Community 44 - "newWatchPattern"
Cohesion: 0.21
Nodes (17): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+9 more)

### Community 45 - "update_test.go"
Cohesion: 0.06
Nodes (46): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), assetsFor(), checkUpdatePreconditions(), compareStableVersions() (+38 more)

### Community 46 - "Read"
Cohesion: 0.08
Nodes (22): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+14 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "watchSet"
Cohesion: 0.19
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors

### Community 49 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "codeedit.go"
Cohesion: 0.13
Nodes (27): class, diffModel, editorModel, classOf(), dedent(), Edit(), leadingIndent(), lines() (+19 more)

### Community 52 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (28): listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), listenForGuardEvents(), tickGuardStats(), syncViewport(), formatAddr(), formatNetProto() (+20 more)

### Community 55 - "TrustGuard"
Cohesion: 0.06
Nodes (16): tempGrant, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows() (+8 more)

### Community 56 - "unit.go"
Cohesion: 0.09
Nodes (29): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+21 more)

### Community 57 - "absPath"
Cohesion: 0.16
Nodes (3): IntegrationSuite, guardDeltaEvents(), absPath()

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (51): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), TestSetInspectorsBounded(), appendInodeRow() (+43 more)

### Community 59 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.28
Nodes (4): column(), repeatKey(), reposition(), Finder

### Community 60 - "NetworkMonitor"
Cohesion: 0.17
Nodes (6): NetEventTypes(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 61 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "GuardEvent"
Cohesion: 0.06
Nodes (22): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+14 more)

### Community 64 - "expandPlaceholders"
Cohesion: 0.18
Nodes (7): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 65 - "pinstate.go"
Cohesion: 0.07
Nodes (44): configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder() (+36 more)

### Community 66 - "gui/session.go"
Cohesion: 0.05
Nodes (29): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+21 more)

### Community 67 - "inode_dev"
Cohesion: 0.18
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path() (+7 more)

### Community 68 - "Guard"
Cohesion: 0.07
Nodes (15): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), exeRow, rootHandle, openRoot(), TestWithChmodDropWriteNeedsReadOnly() (+7 more)

### Community 69 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - ".judgeProcess"
Cohesion: 0.14
Nodes (13): codeVouched(), isTraced(), judgePasses(), loaderEnv(), parseProcs(), roFor(), TestCodeVouched(), TestJudgePasses() (+5 more)

### Community 73 - "globBuilder"
Cohesion: 0.22
Nodes (8): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestGlobReservationsLibTrusted(), TestParseGlobName()

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "RefreshSection"
Cohesion: 0.09
Nodes (26): LibraryBlock, RefreshOptions, SectionScan, findSectionEnd(), isLibDirectiveText(), ParseSectionHeaderPath(), quotePath(), unquotePath() (+18 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (43): AddServeFlags(), CheckBPFLSM(), CheckEBPF(), credentialFlagState(), isLoopbackHost(), MakeDisplayPaths(), ParseEventsFlag(), ParseServeFlags() (+35 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 81 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 82 - "CandidateDir"
Cohesion: 0.17
Nodes (11): bunEntryConfigured(), CandidateDir, TestExtraWatchPathsForSkipsMissing(), TestPathsFor(), catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath() (+3 more)

### Community 83 - "wizard.go"
Cohesion: 0.06
Nodes (41): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), decryptDirectories(), decryptStep() (+33 more)

### Community 85 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.26
Nodes (10): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle() (+2 more)

### Community 86 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 87 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (20): tempMaskOp, tempRow, SupersededBinary, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants() (+12 more)

### Community 88 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "catalogRefresher"
Cohesion: 0.11
Nodes (15): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps() (+7 more)

### Community 92 - "time.Time"
Cohesion: 0.07
Nodes (42): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), runBPFCheck(), runDaemon() (+34 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.14
Nodes (8): foreignPinnedLink(), NewMountVouch(), eventsetSummary(), eventTypeKey(), modeLabel(), NewNetGuard(), hook, NetGuard

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 100 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "time.Duration"
Cohesion: 0.27
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "runNetworkGuard"
Cohesion: 0.10
Nodes (18): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode (+10 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "MonitorUseCase"
Cohesion: 0.20
Nodes (5): runHeadless(), MonitorRepository, MonitorUseCase, NewMonitorUseCase(), TestMonitorUseCaseStartError()

### Community 110 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.22
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 111 - "KernelDev"
Cohesion: 0.21
Nodes (11): fdStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices() (+3 more)

### Community 112 - "find.go"
Cohesion: 0.24
Nodes (8): findButton, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "EventType"
Cohesion: 0.12
Nodes (11): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), Run(), DecodeBpfEvent(), EventTypes() (+3 more)

### Community 116 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "markReady"
Cohesion: 0.25
Nodes (9): markReady(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows() (+1 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "mountVouchesOwnership"
Cohesion: 0.29
Nodes (7): hostView(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "guardLSMHooks"
Cohesion: 0.33
Nodes (3): trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 144 - "mountTable"
Cohesion: 0.24
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 145 - "buildOneGuard"
Cohesion: 0.07
Nodes (45): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded() (+37 more)

### Community 146 - "ldConfParser"
Cohesion: 0.50
Nodes (3): ldConfParser, ldConfInclude(), ldSoConfDirs()

### Community 147 - "render.go"
Cohesion: 0.12
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 148 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "filevault.go"
Cohesion: 0.16
Nodes (23): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey() (+15 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "integrationtests/guard_test.go"
Cohesion: 0.38
Nodes (5): attrOpCase, guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 154 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 155 - "io.Reader"
Cohesion: 0.17
Nodes (11): checkBtrfsKeys(), BtrfsMounted(), FirstOnBtrfs(), mountinfoDev(), mountinfoHasFstype(), OnBtrfs(), SuperblockDev(), SuperblockDevPath() (+3 more)

### Community 156 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 157 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (57): setConfinedHomes(), confirmUserPlaced(), VettedInode, Multicall, MulticallKind, ExeKey(), IsMulticallRefusal(), multicallFamily() (+49 more)

### Community 161 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 163 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 165 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 168 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 184 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 185 - "SanitizeText"
Cohesion: 0.11
Nodes (13): multicallAdmissible(), ownerPaths(), runHeadless(), runHeadless(), binaryVetter, vetDecision, launchKey(), TestLaunchRulesFitTheKey() (+5 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 455 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `watchSet` to `Daemon configuration`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12096774193548387 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `watchSet`, `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.028603034134007586 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `newWatchPattern` to `watchSet`, `catalogRefresher`, `go_pkg_path_filepath`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05171264827402373 - nodes in this community are weakly interconnected._