# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 436 files · ~1,614,521 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5525 nodes · 23364 edges · 216 communities (179 shown, 37 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2050 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0acc4290`
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
- catalogrefresh.go
- Resource
- .exportTable
- sanitizeTerminalText
- safeio.go
- fscrypt/orphans.go
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- DaemonUseCase
- uninstall_system.go
- NetGuardEvent
- guard.bpf.c
- Ledger
- daemon_multicall_test.go
- inodeChain
- codeedit.go
- newTestTarget
- Serve
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- navItem
- github.com/spf13/cobra.Command
- IntegrationSuite
- BinaryEntry
- controlServer
- fileedit_model_test.go
- Guard
- IntegrationSuite
- absPath
- monitorUnitTest
- IntegrationSuite
- update_test.go
- vaultFS
- stdio
- golang.org/x/sys/unix.Stat_t
- sync.Mutex
- LibraryClosure
- pinstate.go
- update.go
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- install/install.go
- systemd.go
- GuardInodeKey
- .policyTools
- buildOneGuard
- pill
- exe_supersede.h
- trustManager
- installPage
- Vault
- context.Context
- MonitorRepository
- changelog.go
- FileTarget
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- Network
- filevault.go
- newWatchPattern
- NetEventType
- process_vm_readv.c
- GuardEvent
- globBuilder
- NetworkMonitor
- backups.go
- Monitor
- hashFileShared
- Finder
- auth.go
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- netModel
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- headless.go
- policyEditor
- session_test.go
- Well-known bypass classes
- journal
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- find.go
- Build & sign release assets (reusable workflow)
- icon
- KernelDev
- mountTable
- NetworkGuardRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- GlobReservations
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- Highlighter
- filevault_test.go
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- CandidateDir
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- .decryptFileWithProgress
- trace-app-libs.sh
- IntegrationSuite
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- io.Reader
- tempGrantSetup
- pathCache
- EventFanout
- net_tester/main.go
- .eventTypeAllowed
- activity
- wipe
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- layout.go
- downloadAndVerify
- IntegrationSuite
- Installation
- os.File
- graphify-refresh.sh
- time.Duration
- parseGuardEvent
- TestMain
- WithSuspectsFrom
- feed.go
- shell
- IntegrationSuite
- integrationtests/networkguard_test.go
- BinaryStat
- .refreshOnce
- IntegrationSuite
- Limitations
- Troubleshooting
- RawTarget
- mountsUnitTest
- Daemon configuration
- edit-protected
- find_test.go
- configedit_test.go
- Compatibility
- How it works
- writeWithin
- launchKey
- bleedLayout
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- changelog_test.go
- NewDumpHook
- markReady
- monitor: see what touches a path
- kernel-prof.sh
- Execute
- github.com/stretchr/testify/suite.Suite
- IntegrationSuite
- wm.go
- mustSubkey
- .readEvent
- wm_test.go
- Desktop GUI
- testLimiter
- lockAndDeprovision
- resolveBits
- bpfLsmListed
- render.go
- TestBatchInodeRows
- harnessSuite

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `shell` - 66 edges
10. `fileEditModel` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (216 total, 37 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (23): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, multicallCacheKey (+15 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (96): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), discordCatalogEntry(), groupedDiscordConf() (+88 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "confgen.go"
Cohesion: 0.05
Nodes (76): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), collectDiffAdditions(), runDiffCatalog(), TestDiffMergePreservesExistingAndParses(), selectAndEditConfig(), addManualDirectories() (+68 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (82): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.08
Nodes (66): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), configPaths(), pathCovered() (+58 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (23): addErrKind, bpfGuardEvent, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf() (+15 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.26
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "catalogrefresh.go"
Cohesion: 0.34
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 11 - "Resource"
Cohesion: 0.07
Nodes (28): newBinaryVetter(), openBinaryVetter(), ownerPaths(), hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries() (+20 more)

### Community 12 - ".exportTable"
Cohesion: 0.11
Nodes (12): installUI, multiAnswer, scrollable(), encodeTable(), TestEncodeTable(), ask(), diffText(), orDefault() (+4 more)

### Community 13 - "sanitizeTerminalText"
Cohesion: 0.11
Nodes (14): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), formatType() (+6 more)

### Community 14 - "safeio.go"
Cohesion: 0.07
Nodes (40): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+32 more)

### Community 15 - "fscrypt/orphans.go"
Cohesion: 0.23
Nodes (18): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies(), orphanedProtectors() (+10 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "DaemonUseCase"
Cohesion: 0.07
Nodes (9): GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "uninstall_system.go"
Cohesion: 0.11
Nodes (25): TestDaemonUnitWithMetadataOutput(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents() (+17 more)

### Community 24 - "NetGuardEvent"
Cohesion: 0.09
Nodes (21): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNetGuard (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "Ledger"
Cohesion: 0.07
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+9 more)

### Community 27 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 28 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 29 - "codeedit.go"
Cohesion: 0.19
Nodes (26): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+18 more)

### Community 30 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.18
Nodes (10): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), setGutter(), stylesOf() (+2 more)

### Community 34 - "navItem"
Cohesion: 0.08
Nodes (5): navItem, navItemRenderer, asset(), coloredIcon(), newNavItem()

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (38): AddServeFlags(), CheckEBPF(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseEventsFlag(), ParseServeFlags(), requestedBoolFlag() (+30 more)

### Community 37 - "BinaryEntry"
Cohesion: 0.07
Nodes (25): deferredBinary, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), confinedEntry(), modeString() (+17 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (20): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), startControlServer() (+12 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.13
Nodes (32): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), openRootT(), selectNamed() (+24 more)

### Community 40 - "Guard"
Cohesion: 0.05
Nodes (20): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), exeRow, rootHandle, tempGrant, tempMaskOp (+12 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (3): IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update_test.go"
Cohesion: 0.12
Nodes (12): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanStable() (+4 more)

### Community 46 - "vaultFS"
Cohesion: 0.13
Nodes (10): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries(), TestDefaultNewFileMode(), TestIsBinaryContent(), writeFileInPlace() (+2 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 49 - "sync.Mutex"
Cohesion: 0.07
Nodes (22): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newControlManager(), startControlManager() (+14 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 51 - "pinstate.go"
Cohesion: 0.17
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 52 - "update.go"
Cohesion: 0.15
Nodes (17): showChangelog(), applyUpdate(), checkUpdatePreconditions(), confirmUnit(), confirmUpdate(), InstalledVersion(), isTerminal(), parseChecksum() (+9 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.13
Nodes (43): grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource(), applyAllowLib(), applyConfigLine() (+35 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (16): listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), clamp(), syncViewport(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented() (+8 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (26): trustRows, checkBtrfsKeys(), ensureBtrfsLayout(), deleteResKeys(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust(), supersedeMaps() (+18 more)

### Community 56 - "install/install.go"
Cohesion: 0.10
Nodes (30): confirmMasterKeyOverwrite(), runGenKey(), promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups() (+22 more)

### Community 57 - "systemd.go"
Cohesion: 0.04
Nodes (72): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), preflightDeployedBinary(), confirmRunPrereq(), installConfig(), revertSystemFiles(), deprovisionKind (+64 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (46): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+38 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.09
Nodes (34): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions() (+26 more)

### Community 61 - "pill"
Cohesion: 0.05
Nodes (11): appTheme, binaryRow, iconTile, iconTileRenderer, panelRenderer, pill, pillRenderer, targetRow (+3 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "trustManager"
Cohesion: 0.10
Nodes (21): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+13 more)

### Community 64 - "installPage"
Cohesion: 0.06
Nodes (20): binSource, installAction, installedBinary, installPage, logHook, logLine, runView, stagedBinary (+12 more)

### Community 65 - "Vault"
Cohesion: 0.06
Nodes (49): CheckBPFLSM(), catchLifecycleSignals(), configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart() (+41 more)

### Community 66 - "context.Context"
Cohesion: 0.17
Nodes (10): reachable(), centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect (+2 more)

### Community 67 - "MonitorRepository"
Cohesion: 0.18
Nodes (5): runHeadless(), MonitorRepository, MonitorUseCase, NewMonitorUseCase(), TestMonitorUseCaseStartError()

### Community 68 - "changelog.go"
Cohesion: 0.29
Nodes (14): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), latestTime(), newerThanInstalled(), newerThanStable() (+6 more)

### Community 69 - "FileTarget"
Cohesion: 0.12
Nodes (12): guardPill(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget(), NormalizeSpec() (+4 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 73 - "time.Time"
Cohesion: 0.18
Nodes (8): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "Network"
Cohesion: 0.11
Nodes (8): binaryList, newBinaryList(), netWarnings(), netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 76 - "filevault.go"
Cohesion: 0.22
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+6 more)

### Community 77 - "newWatchPattern"
Cohesion: 0.18
Nodes (18): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+10 more)

### Community 78 - "NetEventType"
Cohesion: 0.09
Nodes (13): snapNetMonitor, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), NetEventTypes(), ParseNetEventType(), eventsetSummary() (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (19): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuardHeadless(), runGuardTUI(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+11 more)

### Community 81 - "globBuilder"
Cohesion: 0.20
Nodes (7): entryWriters(), symlinkTargetGlob(), underResource(), globBuilder, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 82 - "NetworkMonitor"
Cohesion: 0.16
Nodes (5): NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 83 - "backups.go"
Cohesion: 0.05
Nodes (50): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), restoreBackups(), askDecrypt() (+42 more)

### Community 84 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), newMountTable(), Monitor, watchedInode

### Community 85 - "hashFileShared"
Cohesion: 0.20
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 87 - "auth.go"
Cohesion: 0.05
Nodes (59): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), groupWatchPaths(), relOrPath() (+51 more)

### Community 88 - "highlight_test.go"
Cohesion: 0.16
Nodes (16): fg(), lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "netModel"
Cohesion: 0.13
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.12
Nodes (25): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+17 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): trustHook, foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "headless.go"
Cohesion: 0.25
Nodes (5): sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), Check(), buildTiny(), TestCheck()

### Community 100 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), newNetGuardTab(), optionHint() (+4 more)

### Community 101 - "session_test.go"
Cohesion: 0.15
Nodes (16): invokingIDs(), ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+8 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "journal"
Cohesion: 0.18
Nodes (5): fakeGuard, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "find.go"
Cohesion: 0.20
Nodes (7): findButton, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "icon"
Cohesion: 0.09
Nodes (29): clickPanel, filesPage, handButton, homePage, panel, statusTile, newLogo(), newBinaryRow() (+21 more)

### Community 110 - "KernelDev"
Cohesion: 0.11
Nodes (21): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), KernelDev(), LstatInode() (+13 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetworkGuardRepository"
Cohesion: 0.12
Nodes (15): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "GlobReservations"
Cohesion: 0.17
Nodes (11): reserveChain(), GlobChild, compileGlobNames(), globText(), GlobReservations, ParseGlobName(), resolveChildren(), TestCompileGlobNames_RepeatedPatternOrsBits() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "Highlighter"
Cohesion: 0.12
Nodes (6): diffModel, editorModel, Highlighter, newDiffModel(), newEditorModel(), changelogModel

### Community 121 - "filevault_test.go"
Cohesion: 0.29
Nodes (18): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestStageRecoverySealsContentAndStaysInPlace(), TestUnlockFileInPlaceRecoversFromInterruptedTransform() (+10 more)

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

### Community 127 - "CandidateDir"
Cohesion: 0.06
Nodes (40): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), confSafeMatch(), Discover(), DiscoverSystem(), expandPlaceholders(), BinaryRule (+32 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - ".decryptFileWithProgress"
Cohesion: 0.17
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), TestCopyXattrsPreservesUserAttributes(), copyXattrs()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (17): SupersededBinary, TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+9 more)

### Community 143 - "io.Reader"
Cohesion: 0.20
Nodes (7): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 144 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "activity"
Cohesion: 0.17
Nodes (9): newActivity(), Process, newExeResolver(), TestRingKeepsNewest(), NewNetwork(), Counts, newRing(), activity (+1 more)

### Community 150 - "wipe"
Cohesion: 0.20
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+2 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.16
Nodes (30): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+22 more)

### Community 154 - "EventType"
Cohesion: 0.10
Nodes (11): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+3 more)

### Community 155 - "layout.go"
Cohesion: 0.25
Nodes (8): adaptiveGrid, checkGrid, gapVBox, cellMin(), newCheckGrid(), newGapVBox(), rowHeights(), visible()

### Community 156 - "downloadAndVerify"
Cohesion: 0.23
Nodes (14): assetsFor(), downloadAndVerify(), downloadReleaseFiles(), FetchLatest(), fetchRelease(), FetchVerified(), filterChannel(), resolveAssets() (+6 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (66): multicallAdmissible(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir() (+58 more)

### Community 161 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 162 - "parseGuardEvent"
Cohesion: 0.29
Nodes (8): decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 163 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 164 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (42): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+34 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (25): pageDef, pageID, resizeWatch, shell, sizedDialog, snapEngines, snapNet, newShell() (+17 more)

### Community 168 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 169 - "BinaryStat"
Cohesion: 0.23
Nodes (11): BinaryStat, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches() (+3 more)

### Community 170 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 179 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

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

### Community 194 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 195 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 196 - "markReady"
Cohesion: 0.50
Nodes (4): markReady(), notifySystemdReady(), InodeOccupancy(), WarnInodeOccupancy()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "Execute"
Cohesion: 0.20
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 200 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 202 - "wm.go"
Cohesion: 0.38
Nodes (7): detectWM(), Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), Kind

### Community 203 - "mustSubkey"
Cohesion: 0.33
Nodes (8): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat(), TestSealOpenFileVaultRoundTrip()

### Community 204 - ".readEvent"
Cohesion: 0.25
Nodes (6): NetBpfEvent, Cstr(), FormatAddr(), Ntohs(), reasonText(), TestReasonText()

### Community 205 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 206 - "Desktop GUI"
Cohesion: 0.22
Nodes (9): Actions, Desktop GUI, Files, Home, Install, Network, Tables, The daemon binary (+1 more)

### Community 207 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 208 - "lockAndDeprovision"
Cohesion: 0.32
Nodes (4): applyRawKeyPolicy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 209 - "resolveBits"
Cohesion: 0.43
Nodes (4): statFunc, resolveBits(), statKey(), UpdaterPlan

### Community 210 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 212 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

## Knowledge Gaps
- **32 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 530 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **37 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `newWatchPattern` to `Daemon configuration`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10402941069413606 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `newWatchPattern`, `README.md`?**
  _High betweenness centrality (0.063) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02947845804988662 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `newWatchPattern` to `time.Time`, `sync.Mutex`, `catalogrefresh.go`?**
  _High betweenness centrality (0.062) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._