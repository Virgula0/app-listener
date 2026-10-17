# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 379 files · ~1,572,454 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4763 nodes · 19989 edges · 182 communities (151 shown, 31 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1670 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2f7f2126`
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
- globreserve_test.go
- .render
- searchDirs
- TestMonitorUseCaseLifecycle
- update.go
- highlight_test.go
- SampleContent
- .startGuardStd
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- Hash
- runMonitor
- guard.bpf.c
- runNetworkMonitor
- filevault.go
- inodeChain
- .FilterExistingWhitelist
- install.go
- serve.go
- bpfstats/main.go
- catalogRefresher
- fileedit_model_test.go
- github.com/spf13/cobra.Command
- Pulse
- User
- controlServer
- io.Reader
- serveWebSocket
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- IntegrationSuite
- update_test.go
- codeedit.go
- stdio
- KernelDev
- daemon_multicall_test.go
- planUpdaters
- integrationtests/networkguard_test.go
- netModel
- watchGroup
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- Ledger
- offline.go
- engine
- absPath
- buildOneGuard
- compileGlobNames
- exe_supersede.h
- GuardEvent
- buildTrustedSet
- pinstate.go
- Run
- bufio.Reader
- changelog.go
- hashFileShared
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- watchSet
- btrfsLayoutFrom
- NewEventFanout
- ResolveBtrfsLayout
- IntegrationSuite
- NetworkMonitor
- process_vm_readv.c
- runLiveEdit
- IntegrationSuite
- Read
- wizard.go
- Monitor
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- binaryVetter
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
- writeWithin
- fakeDaemon
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
- render.go
- app-listener Terminal Demo Recording
- IntegrationSuite
- parseResize
- Vault
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- parseGuardEvent
- make test-integration (rootful Docker suite)
- supersede_probe.c
- time.Duration
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- Config
- pathCache
- Guard
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- fscrypt.go
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- find.go
- tempGrantSetup
- Installation
- GuardInodeKey
- graphify-refresh.sh
- RawTarget
- DiscoverInfraBinaries
- integrationtests/guard_test.go
- find_test.go
- bpfLsmListed
- launchKey
- github.com/stretchr/testify/suite.Suite
- .exeKeys
- TestMain
- Limitations
- Troubleshooting
- mountsUnitTest
- Daemon configuration
- edit-protected
- runNetworkGuard
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

## Communities (182 total, 31 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (23): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+15 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (110): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+102 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "confgen.go"
Cohesion: 0.08
Nodes (43): groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), LibraryBlock, Section, findSectionEnd(), findSectionStart(), firstContentLine(), GenerateConf() (+35 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+71 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.11
Nodes (44): dupFD(), installSSHAgent(), installSSHAgentEnv(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), upsertUnitAt(), writeAtNoFollow(), rcTarget (+36 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.14
Nodes (38): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+30 more)

### Community 10 - ".render"
Cohesion: 0.15
Nodes (10): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), setGutter(), stylesOf() (+2 more)

### Community 11 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 13 - "update.go"
Cohesion: 0.28
Nodes (5): isAbsoluteCandidate(), TestCatalogSanity(), LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected(), hint

### Community 14 - "highlight_test.go"
Cohesion: 0.20
Nodes (22): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+14 more)

### Community 15 - "SampleContent"
Cohesion: 0.07
Nodes (38): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), TestDaemonUnitWithMetadataOutput() (+30 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (9): clipLabel(), humanSize(), listEntries(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "Hash"
Cohesion: 0.13
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "runMonitor"
Cohesion: 0.11
Nodes (14): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule(), runHeadless(), runMonitor() (+6 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "runNetworkMonitor"
Cohesion: 0.14
Nodes (8): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase

### Community 27 - "filevault.go"
Cohesion: 0.09
Nodes (44): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+36 more)

### Community 28 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 29 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 30 - "install.go"
Cohesion: 0.05
Nodes (75): reloadDaemonForPasswordChange(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled() (+67 more)

### Community 31 - "serve.go"
Cohesion: 0.14
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.11
Nodes (26): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection() (+18 more)

### Community 33 - "catalogRefresher"
Cohesion: 0.16
Nodes (8): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), startCatalogRefresh(), catalogRefresher, debouncer

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (28): OpenRegularNoFollow(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+20 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.12
Nodes (23): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+15 more)

### Community 36 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 37 - "User"
Cohesion: 0.07
Nodes (42): entryWriters(), reserveChain(), underResource(), bunEntryConfigured(), bunLaunchersFor(), cleanOrphanedFscrypt(), addManualDirectories(), groupCandidates() (+34 more)

### Community 38 - "controlServer"
Cohesion: 0.15
Nodes (10): swapFile(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "io.Reader"
Cohesion: 0.33
Nodes (5): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 40 - "serveWebSocket"
Cohesion: 0.23
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 44 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update_test.go"
Cohesion: 0.05
Nodes (38): assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases(), parseChecksum(), parsePublicKey() (+30 more)

### Community 46 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "KernelDev"
Cohesion: 0.08
Nodes (34): fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), binaryStatOf(), isFuseType(), parseMajorMinor() (+26 more)

### Community 49 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 52 - "netModel"
Cohesion: 0.18
Nodes (9): syncViewport(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg (+1 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (10): clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), RunFileEditorSession(), RunSession(), noticeModel, servedModel (+2 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (24): tempGrant, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust() (+16 more)

### Community 56 - "Ledger"
Cohesion: 0.07
Nodes (20): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+12 more)

### Community 57 - "offline.go"
Cohesion: 0.07
Nodes (43): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptStep(), pickDirsToDecrypt(), revertBunLaunchers() (+35 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (29): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+21 more)

### Community 59 - "absPath"
Cohesion: 0.16
Nodes (3): IntegrationSuite, guardDeltaEvents(), absPath()

### Community 60 - "buildOneGuard"
Cohesion: 0.09
Nodes (36): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+28 more)

### Community 61 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "GuardEvent"
Cohesion: 0.08
Nodes (17): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), runGuardHeadless(), TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+9 more)

### Community 64 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 65 - "pinstate.go"
Cohesion: 0.11
Nodes (31): configureDaemonLogging(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+23 more)

### Community 66 - "Run"
Cohesion: 0.05
Nodes (28): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+20 more)

### Community 67 - "bufio.Reader"
Cohesion: 0.20
Nodes (11): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+3 more)

### Community 68 - "changelog.go"
Cohesion: 0.11
Nodes (31): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+23 more)

### Community 69 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.09
Nodes (23): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), copyDirContents(), copyEntry(), copyRegularAt() (+15 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.24
Nodes (15): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+7 more)

### Community 73 - "watchSet"
Cohesion: 0.16
Nodes (12): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), fakeInotify, watchDir (+4 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 76 - "ResolveBtrfsLayout"
Cohesion: 0.22
Nodes (9): checkBtrfsKeys(), ensureBtrfsLayout(), FirstOnBtrfs(), OnBtrfs(), ResolveBtrfsLayout(), SuperblockDev(), SuperblockDevPath(), TestResolveBtrfsLayoutHostBTF() (+1 more)

### Community 78 - "NetworkMonitor"
Cohesion: 0.10
Nodes (14): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventTypeKey() (+6 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "runLiveEdit"
Cohesion: 0.15
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree(), TestWithin() (+7 more)

### Community 82 - "Read"
Cohesion: 0.08
Nodes (22): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+14 more)

### Community 83 - "wizard.go"
Cohesion: 0.09
Nodes (31): deletePostBackups(), restoreBackups(), decryptDirectories(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestDecryptDirectoriesEmpty() (+23 more)

### Community 84 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), newMountTable(), Monitor, watchedInode

### Community 87 - "binaryVetter"
Cohesion: 0.07
Nodes (30): CheckEBPF(), newBinaryVetter(), openBinaryVetter(), awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals() (+22 more)

### Community 88 - "netGuardModel"
Cohesion: 0.21
Nodes (6): computeGuardBinaries(), NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 92 - "time.Time"
Cohesion: 0.08
Nodes (35): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+27 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (19): foreignPinnedLink(), NewMountVouch(), reasonText(), codeVouched(), isTraced(), judgePasses(), loaderEnv(), parseProcs() (+11 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

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
Cohesion: 0.10
Nodes (14): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), symlinkTargetGlob(), expandPlaceholders(), BinaryRule (+6 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - ".refreshOnce"
Cohesion: 0.13
Nodes (18): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), grantNewBinaries() (+10 more)

### Community 110 - "makeReloadHandler"
Cohesion: 0.10
Nodes (19): newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), controlManager, inodeBatchFn (+11 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

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
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "time.Duration"
Cohesion: 0.18
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward() (+4 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (18): SupersededBinary, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess (+10 more)

### Community 144 - "Config"
Cohesion: 0.08
Nodes (41): ServeConfig, sameConfig(), newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted() (+33 more)

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 146 - "Guard"
Cohesion: 0.04
Nodes (33): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, rootHandle, tempMaskOp (+25 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "fscrypt.go"
Cohesion: 0.05
Nodes (51): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+43 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 154 - "EventType"
Cohesion: 0.11
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run() (+4 more)

### Community 156 - "find.go"
Cohesion: 0.07
Nodes (16): findButton, span, diffModel, editorModel, findAll(), lowerRunes(), newFindInput(), overlay() (+8 more)

### Community 157 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "GuardInodeKey"
Cohesion: 0.03
Nodes (76): multicallAdmissible(), ownerPaths(), checkFreshDir(), openDirNoFollow(), renewDir(), vetDecision, refusedReplacements, VettedInode (+68 more)

### Community 161 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 162 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

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
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

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
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "runNetworkGuard"
Cohesion: 0.13
Nodes (15): ParseNetEventsFlag(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), eventsetSummary(), Mode (+7 more)

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
Cohesion: 0.23
Nodes (6): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), TestDefaultNewFileMode(), writeFileInPlace(), vaultFS

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 455 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `watchSet` to `Daemon configuration`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11596123175070544 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `watchSet`, `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02702275600505689 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `watchSet` to `globreserve_test.go`, `update.go`, `catalogRefresher`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05171264827402373 - nodes in this community are weakly interconnected._