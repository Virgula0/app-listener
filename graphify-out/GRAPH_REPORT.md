# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 379 files · ~1,572,266 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4762 nodes · 19984 edges · 190 communities (162 shown, 28 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1669 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f212b482`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- confgen.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- render.go
- LibraryClosure
- TestMonitorUseCaseLifecycle
- update.go
- install.go
- system.go
- .startGuardStd
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- FileEvent
- guard.bpf.c
- NetEventType
- filevault_test.go
- inodeChain
- .FilterExistingWhitelist
- systemd.go
- serve.go
- bpfstats/main.go
- catalogRefresher
- fileedit_model_test.go
- github.com/spf13/cobra.Command
- time.Duration
- User
- sync.Mutex
- BinaryEntry
- GuardInodeKey
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- IntegrationSuite
- update_test.go
- codeedit.go
- stdio
- golang.org/x/sys/unix.Stat_t
- daemon_multicall_test.go
- planUpdaters
- inode_dev.h
- netModel
- Config
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- SanitizeText
- fscrypt/orphans.go
- engine
- absPath
- buildOneGuard
- MountVouch
- exe_supersede.h
- GuardEvent
- fakeGuardRepo
- pinstate.go
- gui/session.go
- inode_dev
- changelog.go
- hashFileShared
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- watchSet
- btrfsLayoutFrom
- loadDaemonConfig
- NewDumpHook
- IntegrationSuite
- NetworkMonitor
- process_vm_readv.c
- runNonInteractiveLivePut
- IntegrationSuite
- guardModel
- uninstall_system.go
- Monitor
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- hashCache
- netGuardModel
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- IntegrationSuite
- time.Time
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- model
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- openBunTmp
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- expandPlaceholders
- Build & sign release assets (reusable workflow)
- .refreshOnce
- makeReloadHandler
- mountTable
- launch_rule_at
- app-listener Terminal Demo Recording
- IntegrationSuite
- is_event_type_allowed
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- ResolveTempBinary
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- parseGuardEvent
- make test-integration (rootful Docker suite)
- fakeVault
- supersede_probe.c
- liveSession
- trace-app-libs.sh
- IntegrationSuite
- markReady
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- BinaryStat
- trustManager
- pathCache
- Guard
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- filevault.go
- EventType
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- eventUnitTest
- watchPattern
- highlight_test.go
- tempGrantSetup
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- discoverInfra
- .TestCheckHardlinkByInode
- integrationtests/guard_test.go
- find_test.go
- bpfLsmListed
- launchKey
- github.com/stretchr/testify/suite.Suite
- .exeKeys
- ldConfParser
- TestMain
- Limitations
- Troubleshooting
- mountsUnitTest
- Daemon configuration
- edit-protected
- NetGuardEvent
- Compatibility
- How it works
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

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

## Communities (190 total, 28 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (18): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (87): TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), atomicWriteAt(), descendNoFollow() (+79 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "confgen.go"
Cohesion: 0.06
Nodes (56): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), libraryBlocksFromCandidates(), sectionsFromCandidates(), LibraryBlock, RefreshOptions (+48 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.17
Nodes (32): rcTarget, bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock(), EnsureSSHAgentEnv(), openOwnedRegular(), rcBlock() (+24 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "render.go"
Cohesion: 0.12
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 11 - "LibraryClosure"
Cohesion: 0.15
Nodes (14): defaultLibDirs(), elfInterp(), fileExists(), isScript(), LdSoPreloadPaths(), LibraryClosure(), resolveSoname(), rootOwnedInode() (+6 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.26
Nodes (10): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle() (+2 more)

### Community 13 - "update.go"
Cohesion: 0.32
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 14 - "install.go"
Cohesion: 0.13
Nodes (25): applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly() (+17 more)

### Community 15 - "system.go"
Cohesion: 0.06
Nodes (52): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), runEditConfig(), TestDaemonUnitWithMetadataOutput(), updateCatalogConfig() (+44 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.12
Nodes (9): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent() (+1 more)

### Community 19 - "fileEditModel"
Cohesion: 0.13
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), fakeDaemon, Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.14
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "FileEvent"
Cohesion: 0.12
Nodes (9): runHeadless(), FileEvent, Run(), MonitorRepository, MonitorUseCase, NewMonitorUseCase(), newFakeMonitorRepo(), TestMonitorUseCaseStartError() (+1 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (56): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with() (+48 more)

### Community 26 - "NetEventType"
Cohesion: 0.09
Nodes (13): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), Ntohs(), ParseNetEventType(), reasonText() (+5 more)

### Community 27 - "filevault_test.go"
Cohesion: 0.11
Nodes (37): collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+29 more)

### Community 28 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 29 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 30 - "systemd.go"
Cohesion: 0.05
Nodes (65): confirmMasterKeyOverwrite(), runGenKey(), reloadDaemonForPasswordChange(), confirmRunPrereq(), resolveFilesystemPrereqs(), revertSystemFiles(), applyUpdate(), checkKeyLen() (+57 more)

### Community 31 - "serve.go"
Cohesion: 0.08
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "catalogRefresher"
Cohesion: 0.11
Nodes (13): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), awaitStartupOrSignal(), startCatalogRefresh() (+5 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (26): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+18 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.05
Nodes (58): AddServeFlags(), CheckBPFLSM(), CheckEBPF(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseEventsFlag(), ParseNetEventsFlag() (+50 more)

### Community 36 - "time.Duration"
Cohesion: 0.27
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 37 - "User"
Cohesion: 0.07
Nodes (50): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+42 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.11
Nodes (20): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), startControlServer() (+12 more)

### Community 39 - "BinaryEntry"
Cohesion: 0.07
Nodes (27): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), BinaryEntry, TestWithChmodDropWriteNeedsReadOnly(), checkBtrfsKeys(), BinariesSummary() (+19 more)

### Community 40 - "GuardInodeKey"
Cohesion: 0.06
Nodes (29): refusedReplacements, TestNearestRootClaims(), TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows(), inodeRowsOf(), nearestRootClaims(), Inspector (+21 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 44 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update_test.go"
Cohesion: 0.05
Nodes (35): assetsFor(), checkUpdatePreconditions(), downloadFile(), downloadReleaseFiles(), fetchReleases(), parseChecksum(), parsePublicKey(), readInstalledVersion() (+27 more)

### Community 46 - "codeedit.go"
Cohesion: 0.12
Nodes (29): class, diffModel, editorModel, classOf(), column(), dedent(), Edit(), leadingIndent() (+21 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 49 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "inode_dev.h"
Cohesion: 0.47
Nodes (4): ctx_ptr(), sb_dev(), guard_inode_free(), trust_inode_free()

### Community 52 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 53 - "Config"
Cohesion: 0.08
Nodes (49): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), findResource(), groupWatchPaths(), relOrPath() (+41 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (10): clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), RunFileEditorSession(), RunSession(), noticeModel, servedModel (+2 more)

### Community 55 - "TrustGuard"
Cohesion: 0.09
Nodes (10): tempGrant, trustHook, trustRows, setSupersedeTrust(), liveTrust(), cStr(), TrustGuard, logTrustDenied() (+2 more)

### Community 56 - "SanitizeText"
Cohesion: 0.06
Nodes (29): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), systemPatterns(), systemRule(), grantNewBinaries(), runHeadless(), runHeadless() (+21 more)

### Community 57 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (25): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), pathWithin(), TestPathWithinMatchesPrefixRule() (+17 more)

### Community 59 - "absPath"
Cohesion: 0.16
Nodes (3): IntegrationSuite, guardDeltaEvents(), absPath()

### Community 60 - "buildOneGuard"
Cohesion: 0.08
Nodes (38): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+30 more)

### Community 61 - "MountVouch"
Cohesion: 0.13
Nodes (10): statFunc, deleteInoKeys(), deleteResKeys(), syncU32Map(), resolveBits(), statKey(), syncMap(), deleteKeys() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (20): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+12 more)

### Community 63 - "GuardEvent"
Cohesion: 0.11
Nodes (11): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), runGuardHeadless(), TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+3 more)

### Community 65 - "pinstate.go"
Cohesion: 0.09
Nodes (38): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+30 more)

### Community 66 - "gui/session.go"
Cohesion: 0.05
Nodes (29): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+21 more)

### Community 67 - "inode_dev"
Cohesion: 0.20
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits(), fill_path(), get_current_exe_inode() (+7 more)

### Community 68 - "changelog.go"
Cohesion: 0.10
Nodes (34): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+26 more)

### Community 69 - "hashFileShared"
Cohesion: 0.19
Nodes (10): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+2 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 73 - "watchSet"
Cohesion: 0.23
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 76 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 78 - "NetworkMonitor"
Cohesion: 0.16
Nodes (6): NetEventTypes(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "runNonInteractiveLivePut"
Cohesion: 0.40
Nodes (5): coveredByWatch(), TestWithin(), within(), readPutContent(), runNonInteractiveLivePut()

### Community 82 - "guardModel"
Cohesion: 0.07
Nodes (29): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+21 more)

### Community 83 - "uninstall_system.go"
Cohesion: 0.05
Nodes (59): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), pickUsers(), askDecrypt(), cleanOrphanedMetadata() (+51 more)

### Community 84 - "Monitor"
Cohesion: 0.19
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.19
Nodes (10): findButton, span, reposition(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 87 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 88 - "netGuardModel"
Cohesion: 0.24
Nodes (4): computeGuardBinaries(), protoString(), netGuardEventMsg, netGuardModel

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 92 - "time.Time"
Cohesion: 0.10
Nodes (29): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+21 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (11): foreignPinnedLink(), NewMountVouch(), FirstOnBtrfs(), OnBtrfs(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 100 - "KernelDev"
Cohesion: 0.11
Nodes (24): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+16 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "expandPlaceholders"
Cohesion: 0.08
Nodes (20): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames() (+12 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - ".refreshOnce"
Cohesion: 0.21
Nodes (11): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines() (+3 more)

### Community 110 - "makeReloadHandler"
Cohesion: 0.16
Nodes (10): newControlManager(), startControlManager(), makeReloadHandler(), noteReloadedSettings(), controlManager, CleanupStalePins(), genOfPinFile(), TestCleanupStalePins() (+2 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "ResolveTempBinary"
Cohesion: 0.20
Nodes (8): confirmUserPlaced(), binaryStatOf(), loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows(), TestResolveTempBinaryClassesRuntime(), TestGrantTemporaryAccessStreamsOnlyGrantedExeEvents()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "liveSession"
Cohesion: 0.16
Nodes (15): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), putConfig(), chooseResources(), eventsLabel(), forwardBinaries() (+7 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "markReady"
Cohesion: 0.25
Nodes (9): markReady(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (15): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+7 more)

### Community 143 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 144 - "trustManager"
Cohesion: 0.11
Nodes (21): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+13 more)

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (17): deferredBinary, exeRow, rootHandle, tempMaskOp, tempRow, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+9 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "filevault.go"
Cohesion: 0.07
Nodes (37): deprovisionKind, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), isFileVaultCiphertext() (+29 more)

### Community 151 - "EventType"
Cohesion: 0.25
Nodes (5): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), EventTypes(), ParseEventType()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "eventUnitTest"
Cohesion: 0.32
Nodes (3): BpfEvent, eventUnitTest, DecodeBpfEvent()

### Community 155 - "watchPattern"
Cohesion: 0.52
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 156 - "highlight_test.go"
Cohesion: 0.15
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 157 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (59): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), VettedInode, Multicall, MulticallKind, ExeKey() (+51 more)

### Community 161 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 162 - "discoverInfra"
Cohesion: 0.29
Nodes (7): discoverInfra(), DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 164 - "integrationtests/guard_test.go"
Cohesion: 0.38
Nodes (5): attrOpCase, guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 165 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 166 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 168 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.16
Nodes (4): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 170 - "ldConfParser"
Cohesion: 0.50
Nodes (3): ldConfParser, ldConfInclude(), ldSoConfDirs()

### Community 171 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "NetGuardEvent"
Cohesion: 0.13
Nodes (13): resolveGuardMode(), runTUI(), eventsetSummary(), Mode, NetGuardEvent, modeLabel(), NewNetGuard(), NetworkGuardRepository (+5 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

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

### Community 199 - "vaultFS"
Cohesion: 0.14
Nodes (13): AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries(), openRootT(), TestDefaultNewFileMode(), TestIsBinaryContent() (+5 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 455 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Daemon configuration` connect `Daemon configuration` to `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1247957839262187 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `watchPattern` to `catalogRefresher`, `mkdirs`, `watchSet`, `update.go`, `Daemon configuration`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.030643513789581207 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05171264827402373 - nodes in this community are weakly interconnected._
- **Should `confgen.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06120218579234973 - nodes in this community are weakly interconnected._